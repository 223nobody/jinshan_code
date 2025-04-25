package handlers

import (
	"bytes"
	"encoding/json"
	"fileservice/logger"
	"time"

	"github.com/gin-gonic/gin"
)

// RequestLogger 中间件定义
func RequestLogger(logDir string) gin.HandlerFunc {
	dailyLogger := logger.NewDailyLogger(logDir)

	return func(c *gin.Context) {

		// 创建响应记录器
		recorder := &responseRecorder{
			ResponseWriter: c.Writer,
			body:           bytes.NewBuffer(nil),
		}
		c.Writer = recorder

		// 原始请求处理
		c.Next()

		// 构造响应数据结构
		var response struct {
			Code int         `json:"code"`
			Msg  string      `json:"msg"`
			Data interface{} `json:"data"`
		}

		// 解析响应体
		if recorder.body.Len() > 0 {
			_ = json.Unmarshal(recorder.body.Bytes(), &response)
		}

		// 构造日志条目
		logEntry := map[string]interface{}{
			"timestamp":   time.Now().Format(time.RFC3339Nano),
			"method":      c.Request.Method,
			"path":        c.Request.URL.Path,
			"status_code": c.Writer.Status(),
			"response": map[string]interface{}{
				"code": response.Code,
				"msg":  response.Msg,
				"data": response.Data,
			},
		}

		// 错误处理
		if len(c.Errors) > 0 {
			logEntry["errors"] = c.Errors.Errors()
		}

		// 异步写入日志
		go func() {
			_ = dailyLogger.Log(logEntry)
		}()
	}
}
