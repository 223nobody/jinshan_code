package config

import (
	"fmt"
	"os"
	"time"

	"github.com/joho/godotenv"
)

type AIConfig struct {
	DeepSeekKey string
	TongyiKey   string
	Timeout     time.Duration
}

const (
	SingleSelect int = 1
	MultiSelect  int = 2
	Coding       int = 3
)

type QuestionRequest struct {
	Model    string `json:"model" binding:"omitempty,oneof=deepseek tongyi"`                  // 非必选，默认tongyi
	Language string `json:"language" binding:"omitempty,oneof=go java python javascript c++ css html"` // 非必选，默认go
	Count    int    `json:"count" binding:"omitempty,min=3,max=10"` 
	Type     int    `json:"type" binding:"omitempty,oneof=1 2 3"`
	Keyword  string `json:"keyword" binding:"required"` // 必选参数
}


type QuestionResponse struct {
	Title   string   `json:"title"`
	Answers []string `json:"answers"`
	Rights  []string `json:"rights"`
}

type QuestionResponses struct {
	Questions []QuestionResponse `json:"questions"`
}

type QuestionRequest1 struct {
	Id       int      `json:"id"`   // 题目ID
	Type     int      `json:"type"` // 1-单选题 2-多选题 3-程序题
	Title    string   `json:"title"`
	Language string   `json:"language"` // 语言
	Answers  []string `json:"answers"`  // 所有选项
	Rights   []string `json:"rights"`   // 正确答案
}
type AILog struct {
    AIRes      QuestionResponses `json:"aiRes"`  
    AIReq      QuestionRequest   `json:"aiReq"`
    Status     string            `json:"status"`
    Error      string            `json:"error,omitempty"`
    AIStartTime string           `json:"aiStartTime"`
    AIEndTime   string           `json:"aiEndTime"`
    AICostTime  string           `json:"aiCostTime"`
}

func LoadConfig() (*AIConfig, error) {
	_ = godotenv.Load() // 自动加载.env文件

	cfg := &AIConfig{
		DeepSeekKey: getEnv("DEEPSEEK_API_KEY"),
		TongyiKey:   getEnv("TONGYI_API_KEY"),
		Timeout:     getEnvAsDuration("API_TIMEOUT", 30*time.Second),
	}

	// 关键修改：只要配置了任一API密钥即可
	if cfg.DeepSeekKey == "" && cfg.TongyiKey == "" {
		return nil, fmt.Errorf("至少需要配置一个API密钥（DEEPSEEK_API_KEY或TONGYI_API_KEY）")
	}

	return cfg, nil
}

// 辅助函数
func getEnv(key string) string {
	return os.Getenv(key)
}

func getEnvAsDuration(key string, defaultValue time.Duration) time.Duration {
	if value := os.Getenv(key); value != "" {
		if parsed, err := time.ParseDuration(value); err == nil {
			return parsed
		}
	}
	return defaultValue
}
