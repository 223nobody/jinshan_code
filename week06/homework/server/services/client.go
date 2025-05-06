package services

import (
	"Server/config"
	"context"
	"errors"
)

type AIService interface {
	GenerateQuestion(ctx context.Context, req config.QuestionRequest) (*config.QuestionResponse, error)
}

type AIServiceImpl struct {
	deepseek *DeepSeekClient
	tongyi   *TongyiClient
}

func NewAIService(cfg *config.AIConfig) AIService {
	return &AIServiceImpl{
		deepseek: NewDeepSeekClient(cfg.DeepSeekKey, cfg.Timeout),
		tongyi:   NewTongyiClient(cfg.TongyiKey, cfg.Timeout),
	}
}

func (s *AIServiceImpl) GenerateQuestion(ctx context.Context, req config.QuestionRequest) (*config.QuestionResponse, error) {
	if req.Language == "" {
		req.Language = "go" // 默认语言为go
	}
	if req.Type == 0 { // 因binding中未使用required，空值会转为0
		req.Type = 1 // 默认单选题
	}
	switch req.Model {
	case "deepseek":
		return s.deepseek.Generate(ctx, req)
	case "tongyi":
		return s.tongyi.Generate(ctx, req)
	case "": // 默认使用通义千问
		return s.tongyi.Generate(ctx, req)
	default:
		return nil, errors.New("不支持的AI模型")
	}
}
