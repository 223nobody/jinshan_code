package main

import (
	"log"
	"os"
	"user-manager/handlers"

	"github.com/gin-gonic/gin"
)

type User struct {
	Name  string `json:"name"`
	Email string `json:"email" binding:"required,email"`
	Age   int    `json:"age"`
}

func main() {
	router := gin.Default()

	//检查user.json文件是否存在
	if _, err := os.Stat("user.json"); os.IsNotExist(err) {
		log.Println("user.json文件不存在，创建文件")
		os.Create("user.json")
	}

	router.GET("/users", handlers.GetUsers)
	router.POST("/users", handlers.CreateUser)
	router.PUT("/users", handlers.UpdateUser)
	router.DELETE("/users/:email", handlers.DeleteUser)

	log.Println("服务器启动，监听端口 :8080")
	router.Run(":8080")

}
