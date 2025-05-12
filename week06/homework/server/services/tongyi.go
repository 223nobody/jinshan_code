package services

import (
	"Server/config"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"reflect"
	"sort"
	"strings"
	"time"

	openai "github.com/sashabaranov/go-openai"
)

const tongyiEndpoint = "https://dashscope.aliyuncs.com/compatible-mode/v1"

type TongyiClient struct {
	client  *openai.Client
	timeout time.Duration
}

func NewTongyiClient(apiKey string, timeout time.Duration) *TongyiClient {
	config := openai.DefaultConfig(apiKey)
	config.BaseURL = tongyiEndpoint
	config.APIVersion = "" 

	return &TongyiClient{
		client:  openai.NewClientWithConfig(config),
		timeout: timeout,
	}
}

func buildTongyiPrompt(req config.QuestionRequest) string {
	var builder strings.Builder

	builder.WriteString(fmt.Sprintf("请生成【%d】道关于【%s】的编程题，要求如下：\n", req.Count, req.Keyword))
	builder.WriteString(fmt.Sprintf("- 题目数量：%d道\n", req.Count))
	builder.WriteString(fmt.Sprintf("- 编程语言：%s\n", req.Language))
	builder.WriteString(fmt.Sprintf("- 题目类型：%s\n", getQuestionTypeText1(req.Type)))

	switch req.Type {
	case config.SingleSelect:
		builder.WriteString("- 必须且仅有一个正确答案，答案字母需从A/B/C/D中选择\n")
	case config.MultiSelect:
		builder.WriteString("- 正确答案数量需在2-4个之间\n- 答案字母必须按A、B、C、D顺序排列\n- 禁止出现重复字母\n")
	default:
		builder.WriteString("- 必须且仅有一个正确答案，答案字母需从A/B/C/D中选择\n")
	}

	builder.WriteString("\n请严格遵循与以下样例相同的JSON格式：\n")
	switch req.Type {
	case config.SingleSelect:
		builder.WriteString(`[
			{
				"title": "关于Golang并发的说法哪个正确？",
				"answers": [
					"A: channel只能传递基本数据类型",
					"B: sync.Mutex适用于读多写少场景",
					"C: WaitGroup的Add()必须在goroutine外调用",
					"D: map的并发读写需要加锁"
				],
				"rights": ["D"]  //有且仅有一个正确答案
			},
			{
                "title": "Go语言切片行为相关说法",
                "answers": [
                    "A: 切片可以指向nil切片",
                    "B: 切片的长度和容量可以相同",
                    "C: 对切片进行append操作会创建新数组",
                    "D: 切片支持负数索引"
                ],
                "rights": ["A"]
            }
		]`)
	case config.MultiSelect:
		builder.WriteString(`[
        {
			"title": "下面有关Python列表操作相关说法正确的是？",
			"answers": [
				"A: 列表推导式比for循环效率更高",
				"B: 切片操作会创建新对象",
				"C: append()会直接修改原列表",
				"D: 列表可以作为字典的键"
			],
			"rights": ["A","B"]
		},
        {
            "title": "下面有关Go语言切片操作相关说法正确的是？",
            "answers": [
                "A: make([]int, 3) 创建一个长度为3的切片",
                "B: 切片操作不会改变底层数组的长度",
                "C: 使用[:]可以复制整个切片",
                "D: 对切片进行追加操作会影响原始切片"
            ],
            "rights": ["A","C","D"]
        }
            ]`)
	default:
		builder.WriteString(`[
        {
			"title": "请设计C语言中的DFS算法应该怎么写？",
			"answers": [
				"A: void dfs(int v) { visited[v] = 1; for(int i=0; i<vertices; i++) if(graph[v][i] && !visited[i]) dfs(i); }",
				"B: void dfs(int v) { if(!visited[v]) { visited[v]=1; for(int i=vertices-1; i>=0; i--) if(graph[v][i]) dfs(i); } }",
				"C: void dfs(int v) { visited[v] = true; for(struct node* p=G[v]; p; p=p->next) if(!visited[p->data]) dfs(p->data); }",
				"D: void dfs(int v) { mark[v] = 1; for(int i=0; i<edges[v].count; i++) if(!mark[edges[v].targets[i]]) dfs(i); }"
			],
			"rights": ["B"]
		},
        {
            "title": "请设计C语言中的BFS算法应该怎么写？",
            "answers": [
                "A: void bfs(int start) { Queue q; enqueue(q, start); mark[start] = 1; while(!isEmpty(q)) { int v = dequeue(q); for(int i=0; i<vertices; i++) if(graph[v][i] && !mark[i]) { enqueue(q, i); mark[i] = 1; } } }",
                "B: void bfs(int start) { Queue q; enqueue(q, start); mark[start] = 1; while(!isEmpty(q)) { int v = dequeue(q); for(int i=0; i<vertices; i++) if(graph[v][i] && !mark[i]) { enqueue(q, i); mark[i] = 1; } } }",
                "C: void bfs(int start) { Queue q; enqueue(q, start); mark[start] = 1; while(!isEmpty(q)) { int v = dequeue(q); for(int i=0; i<vertices; i++) if(graph[v][i] && !mark[i]) { enqueue(q, i); mark[i] = 1; } } }",
                "D: void bfs(int start) { Queue q; enqueue(q, start); mark[start] = 1; while(!isEmpty(q)) { int v = dequeue(q); for(int i=0; i<vertices; i++) if(graph[v][i] && !mark[i]) { enqueue(q, i); mark[i] = 1; } } }"
            ],
            "rights": ["A"]
        }   
        ]`)
	}

	builder.WriteString("\n\n必须遵守：\n")
	builder.WriteString("1. 选项前缀必须严格按A/B/C/D顺序生成（示例错误：A→C→B ❌）\n")
	builder.WriteString("2. 多选题答案必须按字母顺序排列（如['A','C'] ✅，['C','A'] ❌）\n")
	builder.WriteString("3. 每个答案字母只能出现一次（出现重复直接视为错误）\n")
	builder.WriteString("4. 编程题的ABCD四个选项必须是纯代码段\n")
	builder.WriteString("5. 严格保证每次生成的题目的标题和选项都不同\n")
	builder.WriteString("6. 生成题目title必须是提问句,以？结尾\n")

	return builder.String()
}

func parseTongyiResponse(content string, req config.QuestionRequest) (*config.QuestionResponses, error) {
	// 预处理：去除可能存在的杂项
    content = strings.TrimSpace(content)
    if !strings.HasPrefix(content, "[") || !strings.HasSuffix(content, "]") {
        return nil, fmt.Errorf("响应内容不是有效的JSON数组")
    }
    var response []config.QuestionResponse
	if err := json.Unmarshal([]byte(content), &response); err != nil {
		return nil, fmt.Errorf("响应解析失败: %w", err)
	}

	if len(response) != req.Count {
		return nil, fmt.Errorf("题目数量错误，预期 %d 道，实际 %d 道", req.Count, len(response))
	}

	switch req.Type {
	case config.SingleSelect, config.MultiSelect:
		for _, question := range response {

			// 验证答案
			seen := make(map[string]bool)
			for _, r := range question.Rights {
				if seen[r] {
					return nil, fmt.Errorf("答案重复：%s", r)
				}
				seen[r] = true
			}

			if req.Type == config.MultiSelect {
				sorted := make([]string, len(question.Rights))
				copy(sorted, question.Rights)
				sort.Strings(sorted)
				if !reflect.DeepEqual(sorted, question.Rights) {
					return nil, fmt.Errorf("答案必须按字母顺序排列，当前顺序：%v", question.Rights)
				}
			}
		}
	default:
		// 编程题校验逻辑（可根据需要补充）
	}

	return &config.QuestionResponses{
        Questions: response,
    }, nil
}

func (c *TongyiClient) Generate(ctx context.Context, req config.QuestionRequest) (*config.QuestionResponses, error) {
	ctx, cancel := context.WithTimeout(ctx, c.timeout)
	defer cancel()

	if req.Keyword == "" {
		return nil, errors.New("关键字不能为空")
	}
	if req.Type != config.Coding && req.Language == "" {
		return nil, errors.New("编程语言必须指定")
	}

	prompt := buildTongyiPrompt(req)

	request := openai.ChatCompletionRequest{
		Model: "qwen-turbo",
		Messages: []openai.ChatCompletionMessage{
			{
				Role:    openai.ChatMessageRoleSystem,
				Content: "你是一个严格遵循格式要求的编程题库生成助手",
			},
			{
				Role:    openai.ChatMessageRoleUser,
				Content: prompt,
			},
		},
		ResponseFormat: &openai.ChatCompletionResponseFormat{Type: "json_object"},
		Temperature:    0.3,
		MaxTokens:      2000,
	}

	maxRetries := 3
	var resp openai.ChatCompletionResponse
	var err error // 修复：修正变量名拼写错误

	for i := 0; i < maxRetries; i++ {
		resp, err = c.client.CreateChatCompletion(ctx, request)
		if err == nil {
			break
		}

		if isRetriableError(err) {
			wait := time.Duration(i+1) * 2 * time.Second
			fmt.Printf("[TONGYI_WARN] 请求失败，%s后重试... 错误：%v\n", wait, err)
			time.Sleep(wait)
			continue
		}
		break
	}

	if err != nil {
		return nil, fmt.Errorf("API请求失败（尝试%d次）：%w", maxRetries, err)
	}

	rawResponse := resp.Choices[0].Message.Content

	return parseTongyiResponse(rawResponse, req)
}

func isRetriableError(err error) bool {
	var apiErr *openai.APIError
	if errors.As(err, &apiErr) {
		return apiErr.HTTPStatusCode == 429 || apiErr.HTTPStatusCode >= 500
	}
	return errors.Is(err, context.DeadlineExceeded)
}

func getQuestionTypeText1(t int) string {
	if t == config.MultiSelect {
		return "多选题"
	} else if t == config.SingleSelect {
		return "单选题"
	}
	return "编程题"
}