package main

import (
	"fileservice/db"
	"fmt"
	"io"
	"mime"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

const (
	uploadDir     = "./uploads" // 文件存储目录
	maxUploadSize = 100 << 20   // 100MB
	downloadURL   = "/files/download"
	previewURL    = "/files/preview"
	DbName        = "file_service.db" // CET4数据库文件名
)

type FileInfo struct {
	OriginalName string    `json:"original_name"`
	UUIDName     string    `json:"uuid_name"`
	Size         int64     `json:"size"`
	MIMEType     string    `json:"mime_type"`
	UploadTime   time.Time `json:"upload_time"`
}

var dbfile *db.Database // Declare dbfile as a global variable

func main() {
	var err error

	// 删除旧数据库文件（如果存在）
	_ = os.Remove(DbName)
	dbfile, err = db.InitDatabase(DbName)
	if err != nil {
		panic(fmt.Sprintf("初始化数据库失败: %v", err))
	}
	defer dbfile.Close() // 主函数退出前关闭连接

	// 初始化上传目录
	if err := os.MkdirAll(uploadDir, 0755); err != nil {
		panic(fmt.Sprintf("创建上传目录失败: %v", err))
	}

	router := gin.Default()

	// 文件上传
	router.POST("/files/upload", uploadHandler)

	// 文件下载
	router.GET(downloadURL+"/:uuid", downloadHandler)
	// 文件预览
	router.GET(previewURL+"/:uuid", previewHandler)

	// 健康检查,验证是否正常连接
	router.GET("/api/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})

	router.Run(":8081")
}

// 文件上传处理
func uploadHandler(c *gin.Context) {
	// 限制请求体大小
	c.Request.Body = http.MaxBytesReader(c.Writer, c.Request.Body, maxUploadSize)

	file, fileHeader, err := c.Request.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "文件获取失败"})
		return
	}
	defer file.Close()

	// 获取文件信息
	fileInfo := FileInfo{
		OriginalName: fileHeader.Filename,
		Size:         fileHeader.Size,
	}

	// 生成UUID文件名
	uuid := uuid.New().String()
	fileInfo.UUIDName = uuid + filepath.Ext(fileHeader.Filename)

	// 获取MIME类型
	buff := make([]byte, 512)
	if _, err = file.Read(buff); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "文件读取失败"})
		return
	}
	fileInfo.MIMEType = http.DetectContentType(buff)

	// 校验文件类型
	if !isAllowedType(fileInfo.MIMEType) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "不支持的文件类型"})
		return
	}

	// 创建目标文件
	dstPath := filepath.Join(uploadDir, fileInfo.UUIDName)
	dst, err := os.Create(dstPath)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "文件创建失败"})
		return
	}
	defer dst.Close()

	// 重置文件指针
	if _, err = file.Seek(0, io.SeekStart); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "文件处理失败"})
		return
	}

	// 保存文件
	if _, err = io.Copy(dst, file); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "文件保存失败"})
		return
	}

	// 这里可以保存文件信息到数据库
	dbfile.Insert(fileInfo.OriginalName, fileInfo.Size, fileInfo.UUIDName, fileInfo.MIMEType)

	c.JSON(http.StatusOK, gin.H{
		"download_url": downloadURL + "/" + fileInfo.UUIDName,
		"preview_url":  previewURL + "/" + fileInfo.UUIDName,
		"file_info":    fileInfo,
	})
}

// 文件下载处理
func downloadHandler(c *gin.Context) {
	serveFile(c, true)
}

// 文件预览处理
func previewHandler(c *gin.Context) {
	serveFile(c, false)
}

// 通用文件服务
func serveFile(c *gin.Context, download bool) {
	uuid := c.Param("uuid")
	if uuid == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的文件ID"})
		return
	}

	filePath := filepath.Join(uploadDir, uuid)

	// 检查文件是否存在
	fileInfo, err := os.Stat(filePath)
	if os.IsNotExist(err) {
		c.JSON(http.StatusNotFound, gin.H{"error": "文件不存在"})
		return
	}

	// 打开文件
	file, err := os.Open(filePath)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "文件打开失败"})
		return
	}
	defer file.Close()

	// 设置响应头
	contentType := mime.TypeByExtension(filepath.Ext(uuid))
	if contentType == "" {
		contentType = "application/octet-stream"
	}
	c.Header("Content-Type", contentType)

	if download {
		c.Header("Content-Disposition", "attachment; filename="+filepath.Base(uuid))
	} else {
		c.Header("Content-Disposition", "inline")
	}

	c.Header("Content-Length", fmt.Sprintf("%d", fileInfo.Size()))
	c.File(filePath)
}

// 允许的文件类型校验
func isAllowedType(mimeType string) bool {
	allowedTypes := map[string]bool{
		"image/jpeg":       true,
		"image/png":        true,
		"application/jss":  true,
		"application/html": true,
		"application/css":  true,
	}
	return allowedTypes[mimeType]
}
