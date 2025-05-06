package services

import (
	"Server/config"
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
func (c *DeepSeekClient) Generate(ctx context.Context, req config.QuestionRequest) (*config.QuestionResponse, error) {
	ctx, cancel := context.WithTimeout(ctx, c.timeout)
	defer cancel()

	// 构建规范化的提示语
	prompt := buildDeepseekPrompt(req)

	resp, err := c.client.CreateChatCompletion(ctx, openai.ChatCompletionRequest{
		Model: "deepseek-chat",
		Messages: []openai.ChatCompletionMessage{
			{
				Role:    openai.ChatMessageRoleSystem,
				Content: "你是一个非常专业的编程题库生成助手，严格遵循用户的格式要求",
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

	return parseDeepseekResponse(resp.Choices[0].Message.Content,req)
}

// 构建符合参数规范的提示语
func buildDeepseekPrompt(req config.QuestionRequest) string {
    var builder strings.Builder

    // 基础要求
    builder.WriteString(fmt.Sprintf("请生成关于【%s】的编程题，要求如下：\n", req.Keyword))
    builder.WriteString(fmt.Sprintf("- 编程语言：%s\n", req.Language))
    builder.WriteString(fmt.Sprintf("- 题目类型：%s\n", getQuestionTypeText(req.Type)))
    
    // 根据类型细化规则
    switch req.Type {
    case config.SingleSelect:
        builder.WriteString("- 必须且仅有一个正确答案，答案字母需从A/B/C/D中选择\n")
    case config.MultiSelect:
        builder.WriteString("- 正确答案数量需在2-4个之间，答案字母必须按A、B、C、D顺序排列且没有重复字母出现\n")
    default: // 程序题
        builder.WriteString("- 不需要生成选项和答案，同时必须将answers和rights设为null\n")
    }

    // 格式示例根据类型动态展示
    builder.WriteString("\n请严格遵循以下JSON格式：\n")
    switch req.Type {
    case config.SingleSelect:
        builder.WriteString(`{
  "title": "题目内容",
  "answers": ["A:选项1", "B:选项2", "C:选项3", "D:选项4"],
  "rights": ["A"]  //示例错误："A" ❌  ["A","B"] ❌  
}`)
case config.MultiSelect:
    builder.WriteString(`{
"title": "题目内容",
"answers": ["A:选项1", "B:选项2", "C:选项3", "D:选项4"],
"rights": ["A", "C"]  //❗️必须满足：
                    //1.字母严格按A/B/C/D顺序
                    //2.每个字母只能出现一次
}`)
    default: // 程序题
        builder.WriteString(`{
  "title": "程序题内容",
  "answers": null,
  "rights": null  
}`)
    }

   // 强制约束部分增加强调
   builder.WriteString("\n\n❗❗必须遵守：\n")
   builder.WriteString("1. 多选题答案必须按A、B、C、D顺序排列(示例错误:[\"C\", \"A\"]  ❌ 正确：[\"A\",\"C\"] ✅)\n")
   builder.WriteString("2. 单选题必须只能有一个答案(示例错误:[\"A\", \"B\",\"C\", \"D\"]  ❌ 正确：[\"C\"] ✅)\n")
   builder.WriteString("3. 每个答案字母必须唯一（出现重复视为严重错误）\n")
   builder.WriteString("4. 选项前缀必须严格按A/B/C/D顺序生成\n")
   builder.WriteString(fmt.Sprintf("5. 题目选项内容需与【%s】紧密相关\n", req.Keyword))

    return builder.String()
}
// 修改函数签名传递请求类型
func parseDeepseekResponse(content string, req config.QuestionRequest) (*config.QuestionResponse, error) {
    var response config.QuestionResponse
    if err := json.Unmarshal([]byte(content), &response); err != nil {
        return nil, fmt.Errorf("响应解析失败: %w", err)
    }

    // 根据题目类型验证
    switch req.Type {
    case config.SingleSelect, config.MultiSelect:
        // 选项数量校验
        if len(response.Answers) != 4 {
            return nil, fmt.Errorf("必须提供4个选项，实际收到%d个", len(response.Answers))
        }

        // 选项前缀验证
        for i := 0; i < 4; i++ {
            expected := fmt.Sprintf("%s:", string(rune('A'+i)))
            if !strings.HasPrefix(response.Answers[i], expected) {
                return nil, fmt.Errorf("选项%d格式错误，应以'%s'开头", i+1, expected)
            }
        }

        // 答案合法性
        validAnswers := map[string]bool{"A": true, "B": true, "C": true, "D": true}
        seen := make(map[string]bool)
        for _, r := range response.Rights {
            if !validAnswers[r] {
                return nil, fmt.Errorf("非法答案标识：%s", r)
            }
            if seen[r] {
                return nil, fmt.Errorf("答案重复：%s", r)
            }
            seen[r] = true
        }

        // 严格校验多选题
    if req.Type == config.MultiSelect {
        // 答案数量检查
        if len(response.Rights) < 2 || len(response.Rights) > 4 {
            return nil, fmt.Errorf("多选题需要2-4个答案，当前数量：%d", len(response.Rights))
        }

        // 字母顺序校验（强化版）
        prevChar := 'A' - 1 // 初始值小于A
        for _, char := range response.Rights {
            current := []rune(char)[0]
            if current <= prevChar {
                return nil, fmt.Errorf("答案顺序错误：%v 要求严格递增", response.Rights)
            }
            prevChar = current
        }
    }

    default: // 程序题
        if response.Answers != nil || response.Rights != nil {
            return nil, fmt.Errorf("编程题必须设置answers和rights为null")
        }
    }

    return &response, nil
}


// getQuestionTypeText 根据题目类型返回对应的文本描述
// 这里假设题目类型1为单选题，2为多选题
func getQuestionTypeText(t int) string {
	if t == config.MultiSelect {
		return "多选题"
	} else if t == config.SingleSelect {
		return "单选题"
	} else {
		return "程序题"
	}
}
