package handlers

import (
	"fileservice/api"
	"fileservice/storage"
	"net/http"

	"github.com/gin-gonic/gin"
)

type StatsHandler struct {
	db *storage.Database
}

// 类型统计响应结构
type typeStatResponse struct {
	Type      string `json:"type"`
	FileCount int    `json:"file_count"`
	TotalSize int64  `json:"total_size"` // 使用int64防止溢出
}

func NewStatsHandler(db *storage.Database) *StatsHandler {
	return &StatsHandler{db: db}
}

func (h *StatsHandler) Summary(c *gin.Context) {
	var result struct {
		TotalFiles int `db:"total_files"`
		TotalSize  int `db:"total_size"`
	}

	err := h.db.Get(&result, `
		SELECT 
			COUNT(*) as total_files,
			COALESCE(SUM(size), 0) as total_size
		FROM files
	`)

	if err != nil {
		api.Error(c, http.StatusInternalServerError, "统计失败")
		return
	}

	api.Success(c, result)
}

// ByType 按文件类型统计
func (h *StatsHandler) ByType(c *gin.Context) {
	var stats []struct {
		MimeType  string `db:"mime_type"`
		Count     int    `db:"file_count"`
		TotalSize int64  `db:"total_size"`
	}

	// 执行分组查询
	err := h.db.Select(&stats, `
		SELECT 
			mime_type,
			COUNT(*) as file_count,
			COALESCE(SUM(size), 0) as total_size
		FROM files
		GROUP BY mime_type
		ORDER BY file_count DESC
	`)

	if err != nil {
		api.Error(c, http.StatusInternalServerError, "类型统计失败")
		return
	}

	// 转换响应格式
	var response []typeStatResponse
	for _, s := range stats {
		response = append(response, typeStatResponse{
			Type:      s.MimeType,
			FileCount: s.Count,
			TotalSize: s.TotalSize,
		})
	}

	api.Success(c, gin.H{
		"stats": response,
	})
}
