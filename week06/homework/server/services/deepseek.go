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

func buildDeepseekPrompt(req config.QuestionRequest) string {
	var builder strings.Builder

	builder.WriteString(fmt.Sprintf("请生成【%d】道关于【%s】的编程题，要求如下：\n", req.Count, req.Keyword))
	builder.WriteString(fmt.Sprintf("- 编程语言：%s\n", req.Language))
	builder.WriteString(fmt.Sprintf("- 题目类型：%s\n", getQuestionTypeText(req.Type)))

	switch req.Type {
	case config.SingleSelect:
		builder.WriteString("- 必须且仅有一个正确答案，答案字母需从A/B/C/D中选择\n")
	case config.MultiSelect:
		builder.WriteString("- 正确答案数量需在2-4个之间，答案字母必须按A、B、C、D顺序排列且没有重复字母出现\n")
	default:
		builder.WriteString("- 不需要生成选项和答案，同时必须将answers和rights设为null\n")
	}

	builder.WriteString("\n请严格遵循以下JSON格式：\n")
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

	builder.WriteString("\n\n❗❗必须遵守：\n")
	builder.WriteString("1. 多选题答案必须按A、B、C、D顺序排列\n")
	builder.WriteString("2. 单选题必须只能有一个答案\n")
	builder.WriteString("3. 答案字母必须唯一\n")
	builder.WriteString("4. 选项前缀严格按顺序生成\n")
	builder.WriteString("5. 保证题目和选项不重复\n")
	builder.WriteString("6. 生成题目title必须是提问句,以？结尾\n")

	return builder.String()
}

func parseDeepseekResponse(content string, req config.QuestionRequest) (*config.QuestionResponses, error) {
	// 预处理：去除可能存在的杂项
	content = strings.TrimSpace(content)
	if !strings.HasPrefix(content, "[") || !strings.HasSuffix(content, "]") {
		return nil, fmt.Errorf("响应内容不是有效的JSON数组")
	}

	var items []config.QuestionResponse
	if err := json.Unmarshal([]byte(content), &items); err != nil {
		fmt.Printf("[DEBUG] 原始错误响应内容: %s\n", content)
		return nil, fmt.Errorf("响应解析失败: %w", err)
	}

	if len(items) != req.Count {
		return nil, fmt.Errorf("题目数量错误，预期 %d 道，实际 %d 道", req.Count, len(items))
	}

	switch req.Type {
	case config.SingleSelect, config.MultiSelect:
		for _, question := range items {
			// 验证选项前缀
			for i := 0; i < 4; i++ {
				expected := fmt.Sprintf("%s:", string(rune('A'+i)))
				if !strings.HasPrefix(question.Answers[i], expected) {
					return nil, fmt.Errorf("选项 %d 前缀错误，应以 '%s' 开头", i+1, expected)
				}
			}

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
	}

	return &config.QuestionResponses{
		Questions: items,
	}, nil
}

func (c *DeepSeekClient) Generate(ctx context.Context, req config.QuestionRequest) (*config.QuestionResponses, error) {
	ctx, cancel := context.WithTimeout(ctx, c.timeout)
	defer cancel()

	if req.Keyword == "" {
		return nil, errors.New("关键字不能为空")
	}
	if req.Type != config.Coding && req.Language == "" {
		return nil, errors.New("编程语言必须指定")
	}
	if req.Count > 10 {
		return nil, errors.New("单次生成题目数量不能超过10道")
	}

	prompt := buildDeepseekPrompt(req)

	request := openai.ChatCompletionRequest{
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
		Temperature:    0.3,
		MaxTokens:      2000,
	}

	maxRetries := 3
	var resp openai.ChatCompletionResponse
	var err error

	for i := 0; i < maxRetries; i++ {
		resp, err = c.client.CreateChatCompletion(ctx, request)
		if err == nil {
			break
		}

		if isRetriableError(err) {
			wait := time.Duration(i+1) * 2 * time.Second
			fmt.Printf("[DEEPSEEK_WARN] 请求失败，%s后重试... 错误：%v\n", wait, err)
			time.Sleep(wait)
			continue
		}
		break
	}

	if err != nil {
		return nil, fmt.Errorf("API请求失败（尝试%d次）：%w", maxRetries, err)
	}

	rawResponse := resp.Choices[0].Message.Content
	return parseDeepseekResponse(rawResponse, req)
}

func getQuestionTypeText(t int) string {
	if t == config.MultiSelect {
		return "多选题"
	} else if t == config.SingleSelect {
		return "单选题"
	}
	return "编程题"
}