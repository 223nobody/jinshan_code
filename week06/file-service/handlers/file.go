package handlers

import (
	"bytes"
	"database/sql"
	"errors"
	"fileservice/api"
	"fileservice/storage"
	"fileservice/utils"
	"fmt"
	"io"
	"mime"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

// 文件列表响应结构
type FileListResponse struct {
	UUID      string    `json:"uuid" db:"uuid"`
	Filename  string    `json:"filename" db:"filename"`
	Size      int64     `json:"size" db:"size"`
	Type      string    `json:"type" db:"type"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
}

type FileHandler struct {
	db        *storage.Database
	fileStore *storage.FileStore
}

type responseRecorder struct {
	gin.ResponseWriter
	body *bytes.Buffer
}

var meta struct {
	Filename string `db:"filename"`
	MimeType string `db:"mime_type"`
	Path     string `db:"uuid"`
}

// 新建FileHandler
func NewFileHandler(db *storage.Database, fileStore *storage.FileStore) *FileHandler {
	return &FileHandler{
		db:        db,
		fileStore: fileStore,
	}
}

// Upload 文件上传（支持多文件）
func (h *FileHandler) Upload(c *gin.Context) {
	// 1. 解析上传请求
	form, err := c.MultipartForm()
	if err != nil {
		api.Error(c, http.StatusBadRequest, "无效的表单请求")
		return
	}

	// 2. 处理每个上传文件
	files := form.File["files"]
	if len(files) == 0 {
		api.Error(c, http.StatusBadRequest, "未上传任何文件")
		return
	}

	// 3. 初始化事务
	tx, err := h.db.Beginx()
	if err != nil {
		api.Error(c, http.StatusInternalServerError, "服务不可用")
		return
	}

	// 4. 处理上传结果
	var uploadedFiles []FileListResponse
	for _, fileHeader := range files {
		// 4.1 打开文件流
		file, err := fileHeader.Open()
		if err != nil {
			tx.Rollback()
			api.Error(c, http.StatusInternalServerError, "文件读取失败")
			return
		}
		defer file.Close()

		// 4.2 验证文件类型
		if !utils.ValidateFileType(fileHeader) {
			tx.Rollback()
			api.Error(c, http.StatusBadRequest,
				fmt.Sprintf("禁止的文件类型: %s", filepath.Ext(fileHeader.Filename)))
			return
		}

		// 4.3 生成文件UUID
		uuid := utils.GenerateUUID()

		// 4.4 保存到存储系统
		if err := h.fileStore.Save(uuid, file); err != nil {
			tx.Rollback()
			api.Error(c, http.StatusInternalServerError, "文件保存失败")
			return
		}

		// 获取文件扩展名
		ext := filepath.Ext(fileHeader.Filename)

		// 获取完整MIME类型
		fullMimeType := mime.TypeByExtension(ext)

		// 提取主类型
		primaryMime := utils.GetPrimaryMIME(fullMimeType)

		// 4.5 写入数据库记录
		result, err := tx.NamedExec(`
			INSERT INTO files (uuid, filename, size, mime_type)
			VALUES (:uuid, :filename, :size, :mime_type)`,
			map[string]interface{}{
				"uuid":      uuid,
				"filename":  fileHeader.Filename,
				"size":      fileHeader.Size,
				"mime_type": primaryMime,
			})

		if err != nil {
			tx.Rollback()
			api.Error(c, http.StatusInternalServerError, "数据库写入失败")
			return
		}

		// 4.6 收集上传结果
		if id, _ := result.LastInsertId(); id > 0 {
			uploadedFiles = append(uploadedFiles, FileListResponse{
				UUID:      uuid,
				Filename:  fileHeader.Filename,
				Size:      fileHeader.Size,
				Type:      mime.TypeByExtension(filepath.Ext(fileHeader.Filename)),
				CreatedAt: time.Now().UTC(),
			})
		}
	}

	// 5. 提交事务
	if err := tx.Commit(); err != nil {
		api.Error(c, http.StatusInternalServerError, "事务提交失败")
		return
	}

	api.Success(c, gin.H{
		"uploaded": len(uploadedFiles),
		"files":    uploadedFiles,
	})
}

// List 文件列表（分页+类型筛选）
func (h *FileHandler) List(c *gin.Context) {
	// 1. 解析分页参数
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	if page < 1 {
		page = 1
	}
	pageSize := 10
	offset := (page - 1) * pageSize

	// 2. 构建查询条件
	fileType := c.Query("type")
	query := `
		SELECT 
			uuid, 
			filename, 
			size, 
			mime_type as type,
			created_at
		FROM files
		WHERE 1=1
	`
	args := make([]interface{}, 0)

	if fileType != "" {
		query += " AND mime_type = ?"
		args = append(args, fileType)
	}

	// 3. 执行分页查询
	var files []FileListResponse
	query += " ORDER BY created_at DESC LIMIT ? OFFSET ?"
	args = append(args, pageSize, offset)

	if err := h.db.Select(&files, query, args...); err != nil {
		api.Error(c, http.StatusInternalServerError, "查询失败")
		return
	}

	// 修改总数查询部分
	var total int
	countQuery := "SELECT COUNT(*) FROM files"
	countArgs := make([]any, 0)

	if fileType != "" {
		countQuery += " WHERE mime_type = ?"
		countArgs = append(countArgs, fileType)
	}

	// 统一参数传递方式
	err := h.db.Get(&total, countQuery, countArgs...)
	if err != nil {
		api.Error(c, http.StatusInternalServerError, "总数查询失败: "+err.Error())
		return
	}

	api.Success(c, gin.H{
		"data":  files,
		"total": total,
		"page":  page,
		"pages": (total + pageSize - 1) / pageSize,
	})
}

// Download 文件下载
func (h *FileHandler) DownloadByUuid(c *gin.Context) {
	uuid := c.Param("uuid")

	// 1. 查询文件元数据
	err := h.db.Get(&meta, `
		SELECT filename, uuid 
		FROM files 
		WHERE uuid = ?
	`, uuid)

	if err != nil {
		api.Error(c, http.StatusNotFound, "文件不存在")
		return
	}

	// 2. 生成下载路径
	downloadDir := "downloads"
	if err := os.MkdirAll(downloadDir, 0755); err != nil {
		api.Error(c, http.StatusInternalServerError, "创建下载目录失败: "+err.Error())
		return
	}
	downloadPath := filepath.Join(downloadDir, meta.Filename)

	// 3. 从上传存储复制到下载目录
	srcFile, err := h.fileStore.Get(meta.Path) // 原始上传文件
	if err != nil {
		api.Error(c, http.StatusInternalServerError, "文件获取失败")
		return
	}
	defer srcFile.Close()

	// 创建目标文件（同名直接覆盖）
	dstFile, err := os.Create(downloadPath)
	if err != nil {
		api.Error(c, http.StatusInternalServerError, "创建下载副本失败")
		return
	}
	defer dstFile.Close()

	// 4. 执行复制
	if _, err := io.Copy(dstFile, srcFile); err != nil {
		api.Error(c, http.StatusInternalServerError, "文件复制失败")
		return
	}

	// 5. 发送文件给客户端
	c.FileAttachment(downloadPath, meta.Filename)
}

// DownloadById 根据文件ID下载文件
func (h *FileHandler) DownloadById(c *gin.Context) {
	// 1. 从请求参数中获取文件ID
	id := c.Param("id")
	if id == "" {
		api.Error(c, http.StatusBadRequest, "文件ID不能为空")
		return
	}
	// 2. 查询文件元数据（通过ID）
	err := h.db.Get(&meta, `
        SELECT filename, uuid 
        FROM files 
        WHERE id = ?     
    `, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			api.Error(c, http.StatusNotFound, "文件不存在")
		} else {
			api.Error(c, http.StatusInternalServerError, "数据库查询失败: "+err.Error())
		}
		return
	}

	// 3. 创建下载目录
	downloadDir := "downloads"
	if err := os.MkdirAll(downloadDir, 0755); err != nil {
		api.Error(c, http.StatusInternalServerError, "创建下载目录失败: "+err.Error())
		return
	}

	// 4. 生成下载路径
	downloadPath := filepath.Join(downloadDir, meta.Filename)

	// 5. 从上传存储复制到下载目录
	srcFile, err := h.fileStore.Get(meta.Path) // 使用UUID获取原始文件
	if err != nil {
		api.Error(c, http.StatusInternalServerError, "文件获取失败: "+err.Error())
		return
	}
	defer srcFile.Close()

	// 6. 创建目标文件
	dstFile, err := os.Create(downloadPath)
	if err != nil {
		api.Error(c, http.StatusInternalServerError, "创建下载副本失败: "+err.Error())
		return
	}
	defer dstFile.Close()

	// 7. 执行复制
	if _, err := io.Copy(dstFile, srcFile); err != nil {
		api.Error(c, http.StatusInternalServerError, "文件复制失败: "+err.Error())
		return
	}

	// 8. 发送文件给客户端
	c.FileAttachment(downloadPath, meta.Filename)
}

// Delete 文件删除
func (h *FileHandler) Delete(c *gin.Context) {
	uuid := c.Param("uuid")

	// 1. 启动事务
	tx, err := h.db.Beginx()
	if err != nil {
		api.Error(c, http.StatusInternalServerError, "服务不可用")
		return
	}

	// 2. 获取文件路径
	var path string
	if err := tx.Get(&path, "SELECT uuid FROM files WHERE uuid = ?", uuid); err != nil {
		tx.Rollback()
		api.Error(c, http.StatusNotFound, "文件不存在")
		return
	}

	// 3. 删除数据库记录
	if _, err := tx.Exec("DELETE FROM files WHERE uuid = ?", uuid); err != nil {
		tx.Rollback()
		api.Error(c, http.StatusInternalServerError, "记录删除失败")
		return
	}

	// 4. 删除uploads下的物理文件
	if err := h.fileStore.Delete(path); err != nil {
		tx.Rollback()
		api.Error(c, http.StatusInternalServerError, "文件删除失败")
		return
	}

	// 5. 删除downloads下的物理文件
	downloadDir := "downloads"
	if err := os.MkdirAll(downloadDir, 0755); err != nil {
		api.Error(c, http.StatusInternalServerError, "创建下载目录失败: "+err.Error())
		return
	}
	downloadPath := filepath.Join(downloadDir, meta.Filename)
	if err := os.Remove(downloadPath); err != nil {
		// 忽略文件不存在的错误
		if !os.IsNotExist(err) {
			tx.Rollback()
			api.Error(c, http.StatusInternalServerError, "下载副本删除失败: "+err.Error())
			return
		}
	}

	// 6. 提交事务
	if err := tx.Commit(); err != nil {
		api.Error(c, http.StatusInternalServerError, "事务提交失败")
		return
	}

	api.Success(c, gin.H{
		"deleted": true,
		"uuid":    uuid,
	})
}

func (h *FileHandler) PreviewByUuid(c *gin.Context) {
	uuid := c.Param("uuid")

	err := h.db.Get(&meta, `
        SELECT filename, mime_type, uuid 
        FROM files 
        WHERE uuid = ?
    `, uuid)

	if err != nil {
		api.Error(c, http.StatusNotFound, "文件不存在")
		return
	}

	// 2. 获取文件内容
	file, err := h.fileStore.Get(meta.Path)
	if err != nil {
		api.Error(c, http.StatusInternalServerError, "文件获取失败: "+err.Error()) // 添加错误详情
		return
	}
	defer file.Close()

	// 3. 设置预览头
	c.Header("Content-Type", meta.MimeType)
	c.Header("Content-Disposition", fmt.Sprintf("inline; filename=%q", meta.Filename))

	// 4. 流式传输文件
	if _, err := io.Copy(c.Writer, file); err != nil {
		api.Error(c, http.StatusInternalServerError, "文件传输失败: "+err.Error())
	}
}

// PreviewById 根据文件ID预览文件
func (h *FileHandler) PreviewById(c *gin.Context) {
	// 1. 获取文件ID参数
	id := c.Param("id")
	if id == "" {
		api.Error(c, http.StatusBadRequest, "文件ID不能为空")
		return
	}

	// 2. 查询文件元数据（通过ID）
	err := h.db.Get(&meta, `
        SELECT filename, mime_type, uuid 
        FROM files 
        WHERE id = ?  -- 核心修改点：WHERE条件改为id字段
    `, id)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			api.Error(c, http.StatusNotFound, "文件不存在")
		} else {
			api.Error(c, http.StatusInternalServerError, "数据库查询失败: "+err.Error())
		}
		return
	}

	// 3. 获取文件内容
	file, err := h.fileStore.Get(meta.Path) // 依然使用uuid作为存储路径
	if err != nil {
		api.Error(c, http.StatusInternalServerError, "文件获取失败: "+err.Error())
		return
	}
	defer file.Close()

	// 4. 设置预览头
	c.Header("Content-Type", meta.MimeType)
	c.Header("Content-Disposition", fmt.Sprintf("inline; filename=%q", meta.Filename))

	// 5. 流式传输文件
	if _, err := io.Copy(c.Writer, file); err != nil {
		api.Error(c, http.StatusInternalServerError, "文件传输失败: "+err.Error())
	}
}

func (r *responseRecorder) Write(b []byte) (int, error) {
	r.body.Write(b)
	return r.ResponseWriter.Write(b)
}
