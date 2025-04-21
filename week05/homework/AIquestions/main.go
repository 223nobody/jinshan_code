package main

import (
	"log"

	"AIquestions/config"
	"AIquestions/controllers"
	"AIquestions/services"
	"AIquestions/storage"

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

	// 创建控制器
	ctrl := controllers.NewController(aiService, jsonStorage)

	// 配置路由
	router := gin.Default()
	router.POST("/api/questions/create", ctrl.GenerateQuestion)

	// 健康检查,验证是否正常连接
	router.GET("/api/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})

	log.Println("服务启动于 :8081")
	if err := router.Run(":8081"); err != nil {
		log.Fatal("服务启动失败:", err)
	}
}
