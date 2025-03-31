package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"strconv"
	"strings"
)

type Task struct {
	Content string `json:"content"`
	Done    bool   `json:"done"`
}

const filename = "tasks.json"

func loadTasks() []Task {
	file, err := os.ReadFile(filename)
	if err != nil {
		return []Task{}
	}
	var tasks []Task
	if err := json.Unmarshal(file, &tasks); err != nil {
		fmt.Println("读取任务数据失败:", err)
		return []Task{}
	}
	return tasks
}

func saveTasks(tasks []Task) {
	data, err := json.MarshalIndent(tasks, "", "  ")
	if err != nil {
		fmt.Println("保存任务失败:", err)
		return
	}
	if err := os.WriteFile(filename, data, 0644); err != nil {
		fmt.Println("写入文件失败:", err)
	}
}

func main() {
	tasks := loadTasks()
	fmt.Println("任务管理工具 (输入 help 查看帮助)")
	reader := bufio.NewReader(os.Stdin)

	for {
		fmt.Println()
		fmt.Print("请输入命令 > ")
		input, _ := reader.ReadString('\n')
		input = strings.TrimSpace(input)
		if input == "" {
			continue
		}

		parts := strings.SplitN(input, " ", 2)
		commond := parts[0]
		args := ""
		if len(parts) > 1 {
			args = parts[1]
		}

		switch commond {
		case "add":
			if args == "" {
				fmt.Println("请输入任务内容")
				continue
			}
			tasks = addTask(tasks, args)
			saveTasks(tasks)
		case "list":
			listTasks(tasks)
		case "done":
			if args == "" {
				fmt.Println("请输入任务编号")
				continue
			}
			tasks = markDone(tasks, args)
			saveTasks(tasks)
		case "delete":
			if args == "" {
				fmt.Println("请输入任务编号")
				continue
			}
			tasks = deleteTask(tasks, args)
			saveTasks(tasks)
		case "exit", "quit":
			fmt.Println("已退出CLI任务管理系统")
			return
		case "help":
			showHelp()
		default:
			fmt.Println("未知命令，请输入 help 查看帮助")
		}
	}
}

func showHelp() {
	fmt.Println(`可用命令：
  add <内容>     添加任务
  list           列出未完成任务
  done <编号>    标记任务完成
  delete <编号>  删除任务
  help           显示帮助信息
  exit/quit      退出程序`)
}

func addTask(tasks []Task, content string) []Task {
	Task := Task{Content: content, Done: false}
	fmt.Println("已添加任务: ", content)
	return append(tasks, Task)
}

func listTasks(tasks []Task) {
	fmt.Println("任务列表:")
	for i, task := range tasks {
		status := " "
		if task.Done {
			status = "✓"
		}
		fmt.Printf("%d. [%s] %s\n", i+1, status, task.Content)
	}
}

func markDone(tasks []Task, numStr string) []Task {
	num, err := strconv.Atoi(numStr)
	if err != nil || num < 1 || num > len(tasks) {
		fmt.Println("无效的任务编号")
		return tasks
	}
	index := num - 1
	tasks[index].Done = true
	fmt.Printf("任务 '%s' 已完成\n", tasks[index].Content)
	return tasks
}

func deleteTask(tasks []Task, numStr string) []Task {
	num, err := strconv.Atoi(numStr)
	if err != nil || num < 1 || num > len(tasks) {
		fmt.Println("无效的任务编号")
		return tasks
	}
	index := num - 1
	deletedContent := tasks[index].Content
	newTasks := append(tasks[:index], tasks[index+1:]...)
	fmt.Printf("已删除任务 '%s'\n", deletedContent)
	return newTasks
}
