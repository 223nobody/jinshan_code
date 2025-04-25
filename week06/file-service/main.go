package main

import (
	"fileservice/config"
	"fileservice/handlers"
	"fileservice/storage"
	"log"
	"os"

	"github.com/gin-gonic/gin"
)

func main() {
	// 初始化配置
	cfg := config.Load()

	// 创建存储目录
	if err := os.MkdirAll(cfg.UploadDir, 0755); err != nil {
		log.Fatal("创建上传目录失败: ", err)
	}

	// 初始化数据库
	db, err := storage.InitDB("file_service.db")
	if err != nil {
		log.Fatal("数据库初始化失败: ", err)
	}
	defer db.Close()

	// 初始化文件存储
	fileStore := storage.NewFileStore(cfg.UploadDir)
	fileHandler := handlers.NewFileHandler(db, fileStore)
	statsHandler := handlers.NewStatsHandler(db)

	// 配置Gin
	router := gin.Default()
	router.Use(
		handlers.RequestLogger(cfg.LogDir),
		gin.Recovery(),
	)

	// 文件操作接口
	fileGroup := router.Group("/files")
	{
		fileGroup.POST("/upload", fileHandler.Upload)
		fileGroup.GET("", fileHandler.List)
		fileGroup.GET("/downloadbyuuid/:uuid", fileHandler.DownloadByUuid)
		fileGroup.GET("/previewbyuuid/:uuid", fileHandler.PreviewByUuid)
		fileGroup.GET("/downloadbyid/:id", fileHandler.DownloadById)
		fileGroup.GET("/previewbyid/:id", fileHandler.PreviewById)
		fileGroup.DELETE("/:uuid", fileHandler.Delete)
	}

	// 统计接口
	statsGroup := router.Group("/stats")
	{
		statsGroup.GET("/summary", statsHandler.Summary)
		statsGroup.GET("/by-type", statsHandler.ByType)
	}

	// 健康检查
	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status": "ok",
		})
	})
	router.Run(":8081")

}
