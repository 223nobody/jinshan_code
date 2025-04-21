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
