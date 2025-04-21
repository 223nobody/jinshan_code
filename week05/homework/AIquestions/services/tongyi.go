package services

import (
	"AIquestions/models"
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	openai "github.com/sashabaranov/go-openai"
)

// 阿里云通义千问API地址
const tongyiEndpoint = "https://dashscope.aliyuncs.com/compatible-mode/v1"

type TongyiClient struct {
	client  *openai.Client
	timeout time.Duration
}

func NewTongyiClient(apiKey string, timeout time.Duration) *TongyiClient {
	config := openai.DefaultConfig(apiKey)
	// 修改API地址并添加阿里云所需Header
	config.BaseURL = tongyiEndpoint
	config.APIVersion = "" // 阿里云不需要版本号

	return &TongyiClient{
		client:  openai.NewClientWithConfig(config),
		timeout: timeout,
	}
}

// Generate 生成题目实现（适配阿里云接口）
func (c *TongyiClient) Generate(ctx context.Context, req models.QuestionRequest) (*models.QuestionResponse, error) {
	ctx, cancel := context.WithTimeout(ctx, c.timeout)
	defer cancel()

	// 构建阿里云专用提示语（格式要求可能不同）
	prompt := buildTongyiPrompt(req)

	resp, err := c.client.CreateChatCompletion(ctx, openai.ChatCompletionRequest{
		Model: "qwen-turbo", // 阿里云指定模型名称
		Messages: []openai.ChatCompletionMessage{
			{
				Role:    openai.ChatMessageRoleSystem,
				Content: "你是一个专业的编程题库生成助手，严格遵循用户的格式要求",
			},
			{
				Role:    openai.ChatMessageRoleUser,
				Content: prompt,
			},
		},
		ResponseFormat: &openai.ChatCompletionResponseFormat{Type: "json_object"},
	})

	if err != nil {
		return nil, fmt.Errorf("通义千问API请求失败: %w", err)
	}

	// 阿里云返回结构需特别处理
	return parseTongyiResponse(resp.Choices[0].Message.Content)
}

// 构建适配阿里云的提示语
func buildTongyiPrompt(req models.QuestionRequest) string {
	var builder strings.Builder

	builder.WriteString(fmt.Sprintf("请生成关于【%s】的编程题，要求：\n", req.Keyword))
	builder.WriteString(fmt.Sprintf("- 编程语言：%s\n", req.Language))
	builder.WriteString(fmt.Sprintf("- 题目类型：%s\n", getQuestionTypeText1(req.Type)))
	builder.WriteString("- 选项数量：4个\n\n")
	// 阿里云对格式要求更严格，需明确说明键名
	builder.WriteString(`请严格使用如下JSON格式返回结果：
{
   "title": "题目内容",
  "answers": ["A: 选项1","B: 选项2","C: 选项3","D: 选项4"(例如"D: 数组切片操作会改变原始数组的内容")],
  "rights": [选项一正确对应"A"，选项二正确对应"B"，选项三正确对应"C"，选项四正确对应"D"]
}`)

	return builder.String()
}

// 解析阿里云返回的JSON（适配可能的格式差异）
func parseTongyiResponse(content string) (*models.QuestionResponse, error) {
	var response models.QuestionResponse
	if err := json.Unmarshal([]byte(content), &response); err != nil {
		return nil, fmt.Errorf("响应解析失败: %w", err)
	}

	// 选项数量验证
	if len(response.Answers) != 4 {
		return nil, fmt.Errorf("选项数量必须为4个，当前收到%d个", len(response.Answers))
	}

	// 答案验证
	for _, ans := range response.Rights {
		if string(ans) != "A" && string(ans) != "B" && string(ans) != "C" && string(ans) != "D" {
			return nil, fmt.Errorf("答案越界(仅允许A、B、C、D)： %s", ans)
		}
	}

	return &response, nil
}

// getQuestionTypeText 根据题目类型返回对应的文本描述
// 这里假设题目类型1为单选题，2为多选题
func getQuestionTypeText1(t int) string {
	if t == models.MultiSelect {
		return "多选题"
	}
	return "单选题"
}
