package models

const (
	SingleSelect int = 1
	MultiSelect  int = 2
)

type QuestionRequest struct {
	Model    string `json:"model" binding:"omitempty,oneof=deepseek tongyi"`                  // 非必选，默认tongyi
	Language string `json:"language" binding:"omitempty,oneof=go java python javascript c++"` // 非必选，默认go
	Type     int    `json:"type" binding:"omitempty,oneof=1 2"`
	Keyword  string `json:"keyword" binding:"required"` // 必选参数
}

type QuestionResponse struct {
	Title   string   `json:"title"`
	Answers []string `json:"answers"`
	Rights  []string `json:"rights"`
}

type AILog struct {
	AIStartTime string           `json:"aiStartTime"`
	AIEndTime   string           `json:"aiEndTime"`
	AICostTime  string           `json:"aiCostTime"`
	Status      string           `json:"status"`
	AIReq       QuestionRequest  `json:"aiReq"`
	AIRes       QuestionResponse `json:"aiRes,omitempty"`
	Error       string           `json:"error,omitempty"`
}
