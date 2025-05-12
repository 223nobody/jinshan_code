package services

import (
	"Server/config"
	"context"
	"errors"
)

type AIService interface {
	GenerateQuestion(ctx context.Context, req config.QuestionRequest) (*config.QuestionResponses, error) 
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

func (s *AIServiceImpl) GenerateQuestion(ctx context.Context, req config.QuestionRequest) (*config.QuestionResponses, error) {
	if req.Language == "" {
		req.Language = "go"
	}
	if req.Type == 0 {
		req.Type = 1 
	}
	if req.Count == 0 {
		req.Count = 3        
	}
	switch req.Model {
	case "deepseek":
		return s.deepseek.Generate(ctx, req)
	case "tongyi":
		return s.tongyi.Generate(ctx, req)
	case "": 
		return s.tongyi.Generate(ctx, req)
	default:
		return nil, errors.New("不支持的AI模型")
	}
}