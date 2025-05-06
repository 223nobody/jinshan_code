package controllers

import (
	"Server/api"
	"Server/storage"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"Server/config"

	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
)

type StatsHandler struct {
	db *storage.Database
}

// 分页请求结构体（新增搜索字段）
type PageRequest struct {
	Page     int    `form:"page,default=1"`
	PageSize int    `form:"pageSize,default=10"`
	Search   string `form:"search"` // 新增搜索参数
}

type PageResult struct {
	Total     int            `json:"total"`
	Questions []questionInfo `json:"questions"`
}

type questionInfo struct {
	ID    int    `json:"id"`
	Title string `json:"title"`
	Type  int    `json:"type"`
}

type deleteRequest struct {
	IDs []int `json:"ids" binding:"required,min=1"`
}

func NewStatsHandler(db *storage.Database) *StatsHandler {
	return &StatsHandler{db: db}
}

// 统一处理带搜索的分页请求
func handlePagination(h *StatsHandler, c *gin.Context, typeCondition string) {
	var req PageRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		api.Error(c, http.StatusBadRequest, "参数格式错误")
		return
	}

	if req.Page < 1 || req.PageSize < 1 || req.PageSize > 100 {
		api.Error(c, http.StatusBadRequest, "分页参数超出范围")
		return
	}

	// 构建动态查询条件
	var conditions []string
	var args []interface{}

	// 添加类型条件
	if typeCondition != "" {
		conditions = append(conditions, typeCondition)
	}

	// 添加搜索条件
	if req.Search != "" {
		conditions = append(conditions, "title LIKE ?")
		args = append(args, "%"+req.Search+"%")
	}

	baseQuery := "FROM questions"
	if len(conditions) > 0 {
		baseQuery += " WHERE " + strings.Join(conditions, " AND ")
	}

	// 获取总数
	countQuery := "SELECT COUNT(*) " + baseQuery
	var total int
	if err := h.db.Get(&total, countQuery, args...); err != nil {
		api.Error(c, http.StatusInternalServerError, "获取总数失败")
		return
	}

	// 分页查询
	offset := (req.Page - 1) * req.PageSize
	dataQuery := fmt.Sprintf(`
		SELECT id, title, type 
		%s 
		ORDER BY id DESC 
		LIMIT ? OFFSET ?`, baseQuery)

	// 添加分页参数
	args = append(args, req.PageSize, offset)

	var questions []questionInfo
	if err := h.db.Select(&questions, dataQuery, args...); err != nil {
		api.Error(c, http.StatusInternalServerError, "获取题目列表失败")
		return
	}

	api.Success(c, gin.H{
		"total":     total,
		"questions": questions,
	})
}

// 各类型接口
func (h *StatsHandler) Summary(c *gin.Context) { handlePagination(h, c, "") }
func (h *StatsHandler) ByType1(c *gin.Context) { handlePagination(h, c, "type = 1") }
func (h *StatsHandler) ByType2(c *gin.Context) { handlePagination(h, c, "type = 2") }
func (h *StatsHandler) ByType3(c *gin.Context) { handlePagination(h, c, "type = 3") }

// 批量删除
func (h *StatsHandler) BatchDelete(c *gin.Context) {
	var req deleteRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		api.Error(c, http.StatusBadRequest, "参数格式错误")
		return
	}

	query, args, err := sqlx.In(`
		DELETE FROM questions 
		WHERE id IN (?)
	`, req.IDs)

	if err != nil {
		api.Error(c, http.StatusInternalServerError, "生成查询失败")
		return
	}

	if err := h.db.Exec(query, args...); err != nil {
		api.Error(c, http.StatusInternalServerError, "删除操作失败")
		return
	}

	api.Success(c, gin.H{
		"deleted_ids": req.IDs,
		"message":     "删除成功",
	})
}

// 自主生成题目
func (h *StatsHandler) GenerateQuestion(c *gin.Context) {

	var req config.QuestionRequest1
	// 1. 参数绑定和验证
	if err := c.ShouldBindJSON(&req); err != nil {
		api.Error(c, http.StatusBadRequest, "参数错误: "+err.Error())
		return
	}
	id, err := h.db.CreateQuestion(&req)
	if err != nil {
		api.Error(c, http.StatusBadRequest, "数据库创建失败")
	}

	// 验证题目类型
	if req.Type < 1 || req.Type > 3 {
		api.Error(c, http.StatusBadRequest, "无效的题目类型")
	}

	// 验证选项数量
	if len(req.Answers) != 4 {
		api.Error(c, http.StatusBadRequest, "必须提供4个选项")
	}

	api.Success(c, gin.H{
		"insert_id": id,
		"message":   "插入成功",
		"data": gin.H{
			"type":    req.Type,
			"title":   req.Title,
			"language": req.Language,
			"answers": req.Answers,
			"rights":  req.Rights,
		},
	})
}

func (h *StatsHandler) UpdateQuestion(c *gin.Context) {
	// 1. 参数绑定和验证
	var req config.QuestionRequest1
	if err := c.ShouldBindJSON(&req); err != nil {
		api.Error(c, http.StatusBadRequest, "参数格式错误: "+err.Error())
		return
	}

	// 2. 验证题目类型（匹配图片中的单选/多选类型）
	if req.Type != 1 && req.Type != 2 { // 1-单选 2-多选
		api.Error(c, http.StatusBadRequest, "无效的题目类型")
		return
	}

	// 3. 验证选项完整性（图片中显示4个完整选项）
	if len(req.Answers) != 4 {
		api.Error(c, http.StatusBadRequest, "必须提供4个选项")
		return
	}

	// 4. 验证答案有效性（匹配图片中的A/B/D多选情况）
	validOptions := map[string]bool{"A": true, "B": true, "C": true, "D": true}
	for _, ans := range req.Rights {
		if !validOptions[ans] {
			api.Error(c, http.StatusBadRequest, "存在无效选项标识")
			return
		}
	}

	// 5. 执行更新操作
	affected, err := h.db.UpdateQuestion(&req)
	if err != nil {
		api.Error(c, http.StatusInternalServerError, "更新失败: "+err.Error())
		return
	}

	// 6. 返回结果（匹配图片中的数据结构）
	api.Success(c, gin.H{
		"affected_rows": affected,
		"updated_data": gin.H{
			"type":    req.Type,
			"title":   req.Title,
			"answers": req.Answers,
			"rights":  req.Rights,
		},
	})
}

func (h *StatsHandler) ById(c *gin.Context) {
	// 1. 获取并验证ID
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil || id <= 0 {
		api.Error(c, http.StatusBadRequest, "无效的题目ID: "+idStr)
		return
	}

	// 2. 定义数据库接收结构体
	type dbQuestion struct {
		ID       int    `db:"id"`
		Type     int    `db:"type"`
		Title    string `db:"title"`
		Language string `db:"language"`
		Answers  string `db:"answers"`
		Rights   string `db:"rights"`
	}

	// 3. 执行查询
	var q dbQuestion
	err = h.db.Get(&q, `
        SELECT id, type, title,  language, answers, rights
        FROM questions 
        WHERE id = ?`,
		id,
	)

	if err != nil {
		api.Error(c, http.StatusNotFound, "题目不存在")
		return
	}

	// 4. 解析JSON字段
	var answers []string
	if err := json.Unmarshal([]byte(q.Answers), &answers); err != nil {
		api.Error(c, http.StatusInternalServerError, "选项解析失败")
		return
	}

	var rights []string
	if err := json.Unmarshal([]byte(q.Rights), &rights); err != nil {
		api.Error(c, http.StatusInternalServerError, "答案解析失败")
		return
	}

	// 5. 返回完整数据
	api.Success(c, gin.H{
		"id":       q.ID,
		"type":     q.Type,
		"title":    q.Title,
		"language": q.Language,
		"answers":  answers,
		"rights":   rights,
	})
}
