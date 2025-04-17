package main

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type Response struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	ReqData interface{} `json:"reqData"`
	Data    interface{} `json:"data"`
}

type ReqParams struct {
	X int `json:"x"`
	Y int `json:"y"`
}

func main() {
	r := gin.Default()

	r.GET("/api/sum", func(c *gin.Context) {
		// 获取原始请求参数
		xStr := c.Query("x")
		yStr := c.Query("y")

		// 校验参数是否缺失
		if xStr == "" || yStr == "" {
			c.JSON(http.StatusOK, Response{
				Code:    1,
				Message: "missing required parameters",
				ReqData: gin.H{"x": xStr, "y": yStr},
				Data:    nil,
			})
			return
		}

		// 转换参数并校验格式
		x, err := strconv.Atoi(xStr)
		if err != nil {
			c.JSON(http.StatusOK, Response{
				Code:    2,
				Message: "invalid parameter x",
				ReqData: gin.H{"x": xStr, "y": yStr},
				Data:    nil,
			})
			return
		}

		y, err := strconv.Atoi(yStr)
		if err != nil {
			c.JSON(http.StatusOK, Response{
				Code:    3,
				Message: "invalid parameter y",
				ReqData: gin.H{"x": xStr, "y": yStr},
				Data:    nil,
			})
			return
		}

		// 成功响应
		c.JSON(http.StatusOK, Response{
			Code:    0,
			Message: "success",
			ReqData: ReqParams{X: x, Y: y},
			Data:    x + y,
		})
	})

	r.Run(":8080")
}
