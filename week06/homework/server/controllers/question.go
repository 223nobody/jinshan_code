package controllers

import (
	"Server/config"
	"Server/services"
	"Server/storage"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

type QuestionController struct {
	service services.AIService
	storage storage.Storage
	db      *storage.Database
}

func NewController(service services.AIService, storage storage.Storage, db *storage.Database) *QuestionController {
	return &QuestionController{
		service: service,
		storage: storage,
		db:      db,
	}
}

func (c *QuestionController) GenerateQuestion(ctx *gin.Context) {
	startTime := time.Now()
	var req config.QuestionRequest

	// 参数绑定
	if err := ctx.ShouldBindJSON(&req); err != nil {
		sendError(ctx, http.StatusBadRequest, "参数错误: "+err.Error())
		return
	}

	// 调用服务
	resp, err := c.service.GenerateQuestion(ctx, req)
	logEntry := buildLog(req, resp, err, startTime)
	if saveErr := c.storage.Save(logEntry); saveErr != nil {
		log.Printf("日志存储失败: %v", saveErr)
	}

	//数据库操作被注释掉，封装在函数AddQuestions中
	
	ctx.JSON(http.StatusOK, gin.H{
		"code":  0,
		"msg":   "",
		"aiRes": resp,
	})
}

func buildLog(req config.QuestionRequest, resp *config.QuestionResponses, err error, start time.Time) config.AILog {
    logEntry := config.AILog{
        AIStartTime: start.Format("2006-01-02 15:04:05"),
        AIEndTime:   time.Now().Format("2006-01-02 15:04:05"),
        AICostTime:  fmt.Sprintf("%.2fs", time.Since(start).Seconds()),
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

// 2. 新增批量添加接口
func (c *QuestionController) AddQuestions(ctx *gin.Context) {
    var questions []struct {
        Type     int    `json:"type"`
        Title    string `json:"title"`
        Language string `json:"language"`
        Answers  []string `json:"answers"`
        Rights   []string `json:"rights"`
    }

    if err := ctx.ShouldBindJSON(&questions); err != nil {
        sendError(ctx, http.StatusBadRequest, "参数错误")
        return
    }

    // 开启事务
    tx, err := c.db.Beginx()
    if err != nil {
        sendError(ctx, http.StatusInternalServerError, "服务不可用")
        return
    }

    // 批量插入
    for _, q := range questions {
        answersJSON, _ := json.Marshal(q.Answers)
        rightsJSON, _ := json.Marshal(q.Rights)
        
        _, err = tx.NamedExec(`
            INSERT INTO questions (type, title, language, answers, rights)
            VALUES (:type, :title, :language, :answers, :rights)`,
            map[string]interface{}{
                "type":     q.Type,
                "title":    q.Title,
                "language": q.Language,
                "answers":  string(answersJSON),
                "rights":   string(rightsJSON),
            })
        
        if err != nil {
            tx.Rollback()
            sendError(ctx, http.StatusInternalServerError, "存储失败")
            return
        }
    }

    if err := tx.Commit(); err != nil {
        sendError(ctx, http.StatusInternalServerError, "事务提交失败")
        return
    }

    ctx.JSON(http.StatusOK, gin.H{
        "code": 0,
        "msg":  "添加成功",
    })
}