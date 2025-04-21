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

const deepseekEndpoint = "https://ai.forestsx.top/v1"

type DeepSeekClient struct {
	client  *openai.Client
	timeout time.Duration
}

func NewDeepSeekClient(apiKey string, timeout time.Duration) *DeepSeekClient {
	config := openai.DefaultConfig(apiKey)
	config.BaseURL = deepseekEndpoint

	return &DeepSeekClient{
		client:  openai.NewClientWithConfig(config),
		timeout: timeout,
	}
}

// Generate 生成题目实现（完全匹配接口）
func (c *DeepSeekClient) Generate(ctx context.Context, req models.QuestionRequest) (*models.QuestionResponse, error) {
	ctx, cancel := context.WithTimeout(ctx, c.timeout)
	defer cancel()

	// 构建规范化的提示语
	prompt := buildDeepseekPrompt(req)

	resp, err := c.client.CreateChatCompletion(ctx, openai.ChatCompletionRequest{
		Model: "deepseek-chat",
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
		return nil, fmt.Errorf("DeepSeek API请求失败: %w", err)
	}

	return parseDeepseekResponse(resp.Choices[0].Message.Content)
}

// 构建符合参数规范的提示语
func buildDeepseekPrompt(req models.QuestionRequest) string {
	var builder strings.Builder

	builder.WriteString(fmt.Sprintf("请生成关于【%s】的编程题，要求：\n", req.Keyword))
	builder.WriteString(fmt.Sprintf("- 编程语言：%s\n", req.Language))
	builder.WriteString(fmt.Sprintf("- 题目类型：%s\n", getQuestionTypeText(req.Type)))
	builder.WriteString("- 选项数量：4个\n\n")
	builder.WriteString("请严格使用如下JSON格式返回结果：\n")
	builder.WriteString(`{
  "title": "题目内容",
  "answers": [A: "选项1", B: "选项2", C: "选项3", D: "选项4"],
  "rights": [正确选项对应索引下标对应的字母(例如0对应A，1对应B，2对应C，3对应D)]
}`)

	return builder.String()
}

// 解析并验证响应
func parseDeepseekResponse(content string) (*models.QuestionResponse, error) {
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
func getQuestionTypeText(t int) string {
	if t == models.MultiSelect {
		return "多选题"
	}
	return "单选题"
}
