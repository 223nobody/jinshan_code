package controllers

import (
	"AIquestions/models"
	"AIquestions/services"
	"AIquestions/storage"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

type QuestionController struct {
	service services.AIService
	storage storage.Storage
}

func NewController(service services.AIService, storage storage.Storage) *QuestionController {
	return &QuestionController{
		service: service,
		storage: storage,
	}
}

func (c *QuestionController) GenerateQuestion(ctx *gin.Context) {
	startTime := time.Now()
	var req models.QuestionRequest

	// 参数绑定
	if err := ctx.ShouldBindJSON(&req); err != nil {
		sendError(ctx, http.StatusBadRequest, "参数错误: "+err.Error())
		return
	}

	// 调用服务
	resp, err := c.service.GenerateQuestion(ctx, req)
	logEntry := buildLog(req, resp, err, startTime)

	// 仅记录成功日志
	if err == nil {
		if saveErr := c.storage.Save(logEntry); saveErr != nil {
			log.Printf("日志存储失败: %v", saveErr)
		}
	}

	// 构造响应
	if err != nil {
		sendError(ctx, http.StatusInternalServerError, err.Error())
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"code":  0,
		"msg":   "",
		"aiRes": resp,
	})
}

func buildLog(req models.QuestionRequest, resp *models.QuestionResponse, err error, start time.Time) models.AILog {
	logEntry := models.AILog{
		AIStartTime: start.Format("2006-01-02 15:04:05"),
		AIEndTime:   time.Now().Format("2006-01-02 15:04:05"),
		AICostTime: fmt.Sprintf("%ds", int(time.Since(start).Seconds())),
		AIReq:       req,
		Status:      "success",
	}

	if err != nil {
		logEntry.Status = "failed"
		logEntry.Error = err.Error()
	} else {
		logEntry.AIRes = *resp
	}
	return logEntry
}

func sendError(ctx *gin.Context, code int, msg string) {
	ctx.JSON(code, gin.H{
		"code": code,
		"msg":  msg,
	})
}
