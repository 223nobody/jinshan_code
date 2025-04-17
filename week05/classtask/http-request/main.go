package main

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

// 定义与目标API对应的数据结构
type OriginalResponse struct {
	ErrNo  int         `json:"err_no"`
	ErrMsg string      `json:"err_msg"`
	Data   interface{} `json:"data"`
}

type ModifiedResponse struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

func main() {
	// 初始化 Gin 引擎
	r := gin.Default()

	// 设置路由
	r.GET("/api/category/list", func(c *gin.Context) {
		// 1. 创建HTTP客户端
		client := &http.Client{
			Timeout: 5 * time.Second, // 设置超时时间
		}

		// 2. 请求目标API
		req, _ := http.NewRequest("GET", "https://juejin.cn/course", nil)
		req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/91.0.4472.124 Safari/537.36")

		resp, err := client.Do(req)
		if err != nil {
			c.JSON(http.StatusBadGateway, gin.H{
				"code":    502,
				"message": "无法连接上游服务",
				"data":    nil,
			})
			return
		}
		defer resp.Body.Close()

		// 3. 检查响应状态码
		if resp.StatusCode != http.StatusOK {
			c.JSON(http.StatusBadGateway, gin.H{
				"code":    resp.StatusCode,
				"message": "上游服务返回错误",
				"data":    nil,
			})
			return
		}

		// 4. 读取响应内容
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"code":    500,
				"message": "响应数据读取失败",
				"data":    nil,
			})
			return
		}

		// 5. 解析JSON数据
		var original OriginalResponse
		if err := json.Unmarshal(body, &original); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"code":    500,
				"message": "数据解析失败",
				"data":    nil,
			})
			return
		}

		// 6. 转换响应格式
		modified := ModifiedResponse{
			Code:    original.ErrNo,
			Message: original.ErrMsg,
			Data:    original.Data,
		}

		// 7. 返回转换后的响应
		c.JSON(http.StatusOK, modified)
	})

	// 启动服务
	r.Run(":8080")
}
