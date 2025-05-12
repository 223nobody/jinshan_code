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
		Temperature:    0.3,
		MaxTokens:      500,
	})

	if err != nil {
		return nil, fmt.Errorf("DeepSeek API请求失败: %w", err)
	}

	return parseDeepseekResponse(resp.Choices[0].Message.Content)
}

// 构建符合参数规范的提示语
func buildDeepseekPrompt(req models.QuestionRequest) string {
	var builder strings.Builder

	builder.WriteString(fmt.Sprintf("请生成关于【%s】的编程题，要求如下：\n", req.Keyword))
	builder.WriteString(fmt.Sprintf("- 编程语言：%s\n", req.Language))
	builder.WriteString(fmt.Sprintf("- 题目类型：%s\n", getQuestionTypeText(req.Type)))

	switch req.Type {
	case models.SingleSelect:
		builder.WriteString("- 必须且仅有一个正确答案，答案字母需从A/B/C/D中选择\n")
	case models.MultiSelect:
		builder.WriteString("- 正确答案数量需在2-4个之间，答案字母必须按A、B、C、D顺序排列且没有重复字母出现\n")
	}

	builder.WriteString("\n请严格遵循以下JSON格式：\n")
	switch req.Type {
	case models.SingleSelect:
		builder.WriteString(`
			{
				"title": "关于Golang并发的说法哪个正确？",
				"answers": [
					"A: channel只能传递基本数据类型",
					"B: sync.Mutex适用于读多写少场景",
					"C: WaitGroup的Add()必须在goroutine外调用",
					"D: map的并发读写需要加锁"
				],
				"rights": ["D"]  //有且仅有一个正确答案
			}`)
	case models.MultiSelect:
		builder.WriteString(`
        {
			"title": "下面有关Python列表操作相关说法正确的是？",
			"answers": [
				"A: 列表推导式比for循环效率更高",
				"B: 切片操作会创建新对象",
				"C: append()会直接修改原列表",
				"D: 列表可以作为字典的键"
			],
			"rights": ["A","B"]
		    }`)
	}

	builder.WriteString("\n\n❗❗必须遵守：\n")
	builder.WriteString("1. 多选题答案必须按A、B、C、D顺序排列\n")
	builder.WriteString("2. 单选题必须只能有一个答案\n")
	builder.WriteString("3. 答案字母必须唯一\n")
	builder.WriteString("4. 选项前缀严格按顺序生成\n")
	builder.WriteString("5. 保证题目和选项不重复\n")

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
