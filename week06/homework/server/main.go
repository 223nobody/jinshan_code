package main

import (
	"log"
	"net/http"
	"os"
	"path/filepath"
	"runtime"

	"Server/config"
	"Server/controllers"
	"Server/services"
	"Server/storage"

	"github.com/gin-gonic/gin"
)

func main() {
	// 加载配置
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatal("配置加载失败:", err)
	}

	// 初始化服务
	aiService := services.NewAIService(cfg)
	jsonStorage := storage.NewJSONStorage()

	// 初始化数据库
	db, err := storage.InitDB("question_service.db")
	if err != nil {
		log.Fatal("数据库初始化失败: ", err)
	}
	defer db.Close()

	// 创建控制器
	ctrl := controllers.NewController(aiService, jsonStorage, db)
	statsHandler := controllers.NewStatsHandler(db)

	// 配置路由
	router := gin.Default()

	// 获取项目根目录（关键修改）
	_, mainPath, _, _ := runtime.Caller(0)
	rootDir := filepath.Dir(filepath.Dir(mainPath))

	// 通用CORS配置（开发和生产通用）
	router.Use(func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "http://localhost:3000")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "*") // 允许所有方法
		c.Writer.Header().Set("Access-Control-Allow-Headers", "*") // 允许所有头
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}
		c.Next()
	})

	// 静态资源托管（自动检测）
	distPath := filepath.Join(rootDir, "client", "dist")
	if _, err := os.Stat(distPath); !os.IsNotExist(err) {
		router.Static("/assets", filepath.Join(distPath, "assets"))
		router.StaticFile("/favicon.ico", filepath.Join(distPath, "favicon.ico"))
		router.NoRoute(func(c *gin.Context) {
			c.File(filepath.Join(distPath, "index.html"))
		})
		log.Printf("已加载静态资源: %s", distPath)
	} else {
		log.Printf("未找到静态资源目录: %s", distPath)
	}

	// API路由组
	questionGroup := router.Group("/api/questions")
	{
		questionGroup.POST("/CreateByAI", ctrl.GenerateQuestion)
		questionGroup.POST("/batch-insert", ctrl.AddQuestions) 
		questionGroup.POST("/CreateByHand", statsHandler.GenerateQuestion)
		questionGroup.DELETE("/batch-delete", statsHandler.BatchDelete)
		questionGroup.PUT("/update", statsHandler.UpdateQuestion)
	}

	statsGroup := router.Group("/api/stats")
	{
		statsGroup.GET("/summary", statsHandler.Summary)
		statsGroup.GET("/bytype1", statsHandler.ByType1)
		statsGroup.GET("/bytype2", statsHandler.ByType2)
		statsGroup.GET("/bytype3", statsHandler.ByType3)
		statsGroup.GET("/byid/:id", statsHandler.ById)
	}

	// 健康检查
	router.GET("/api/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	// 启动服务
	port := "8080"
	log.Printf("服务启动于 :%s", port)
	if err := router.Run(":" + port); err != nil {
		log.Fatal("服务启动失败:", err)
	}
}
