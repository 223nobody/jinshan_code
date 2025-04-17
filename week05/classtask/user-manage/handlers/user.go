// handlers/users.go
package handlers

import (
	"encoding/json"
	"errors"
	"net/http"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
)

type User struct {
	Name  string `json:"name"`
	Email string `json:"email" binding:"required,email"`
	Age   int    `json:"age"`
}

var users []User

// 文件操作相关错误处理
func loadUsers() error {
	file, err := os.ReadFile("user.json")
	if err != nil {
		if os.IsNotExist(err) {
			return nil
		}
		return err
	}

	if len(file) == 0 {
		return nil
	}

	if err := json.Unmarshal(file, &users); err != nil {
		return errors.New("invalid user data format")
	}
	return nil
}

func saveUsers() error {
	data, err := json.MarshalIndent(users, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile("user.json", data, 0644)
}

// 辅助函数
func findUser(email string) (int, *User) {
	for i, user := range users {
		if strings.EqualFold(user.Email, email) {
			return i, &user
		}
	}
	return -1, nil
}

// Handlers
func GetUsers(c *gin.Context) {
	if err := loadUsers(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "无法读取用户数据",
		})
		return
	}

	c.JSON(http.StatusOK, users)
}

func CreateUser(c *gin.Context) {
	var newUser User
	if err := c.ShouldBindJSON(&newUser); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "无效的请求数据: " + err.Error(),
		})
		return
	}

	if err := loadUsers(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "无法加载用户数据",
		})
		return
	}

	if _, existing := findUser(newUser.Email); existing != nil {
		c.JSON(http.StatusConflict, gin.H{
			"error": "邮箱地址已存在",
		})
		return
	}

	users = append(users, newUser)

	if err := saveUsers(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "无法保存用户数据",
		})
		return
	}

	c.JSON(http.StatusCreated, newUser)
}

func UpdateUser(c *gin.Context) {
	var updatedUser User
	if err := c.ShouldBindJSON(&updatedUser); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "无效的请求数据",
		})
		return
	}

	if err := loadUsers(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "无法加载用户数据",
		})
		return
	}

	index, existing := findUser(updatedUser.Email) //根据唯一的email进行判断
	if existing == nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "用户不存在",
		})
		return
	}

	// 保留原始邮箱地址
	users[index] = updatedUser

	if err := saveUsers(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "无法保存用户数据",
		})
		return
	}

	c.JSON(http.StatusOK, users[index])
}

func DeleteUser(c *gin.Context) {
	email := c.Param("email")
	if email == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "需要提供邮箱地址",
		})
		return
	}

	if err := loadUsers(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "无法加载用户数据",
		})
		return
	}

	index, existing := findUser(email)
	if existing == nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "用户不存在",
		})
		return
	}

	users = append(users[:index], users[index+1:]...)

	if err := saveUsers(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "无法保存用户数据",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "用户删除成功",
	})
}
