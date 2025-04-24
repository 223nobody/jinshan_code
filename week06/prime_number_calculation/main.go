package main

import (
	"fmt"
	"os"
	"strconv"
	"sync"
	"time"
)

type Task struct {
	a, b int
}

var count int = 0

func main() {
	starttime := time.Now() // 记录程序启动时间
	printTimeNow()
	// 验证参数数量
	if len(os.Args) < 3 {
		fmt.Println("使用方法: go run main.go <起始值> <结束值>")
		os.Exit(1)
	}

	// 解析参数
	start, err := strconv.Atoi(os.Args[1])
	if err != nil {
		fmt.Printf("无效的起始值: %s\n", os.Args[1])
		os.Exit(2)
	}

	end, err := strconv.Atoi(os.Args[2])
	if err != nil {
		fmt.Printf("无效的结束值: %s\n", os.Args[2])
		os.Exit(3)
	}

	filename := "primes_" + strconv.Itoa(start) + "_" + strconv.Itoa(end) + ".txt"

	var wg sync.WaitGroup
	tasks := splitRange(start, end, 4)

	// 启动所有计算任务
	for _, task := range tasks {
		wg.Add(1)
		go calculate(task.a, task.b, filename, &wg)
	}

	// 等待所有计算完成
	wg.Wait()
	fmt.Println("所有计算任务完成！")
	fmt.Printf("找到的素数个数为：%d\n", count)
	fmt.Printf("程序共耗时: %.2f秒\n", time.Since(starttime).Seconds())
}

func Judge(n int, filename string) {
	file, err := os.OpenFile(filename, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		fmt.Printf("打开文件失败: %v\n", err)
		return
	}
	defer file.Close()
	if n == 1 {
		return
	}
	for i := 2; i*i <= n; i++ {
		if n%i == 0 {
			return
		}
	}
	_, err = file.WriteString(fmt.Sprintf("%d ", n))
	count++
	if err != nil {
		fmt.Printf("写入文件失败: %v\n", err)
		return
	}
}

func calculate(a, b int, filename string, wg *sync.WaitGroup) {
	defer wg.Done()
	for i := a; i <= b; i++ {
		Judge(i, filename)
	}
	// time.Sleep(time.Second)
}

func splitRange(start, end int, nTasks int) []Task {
	var tasks []Task
	total := end - start + 1
	step := total / nTasks
	remainder := total % nTasks

	currentStart := start
	for i := 0; i < nTasks; i++ {
		currentStep := step
		if i < remainder {
			currentStep++
		}

		currentEnd := currentStart + currentStep - 1
		if currentEnd > end { // 确保不超过原始end
			currentEnd = end
		}

		tasks = append(tasks, Task{currentStart, currentEnd})
		currentStart = currentEnd + 1

		if currentStart > end { // 提前结束循环
			break
		}
	}
	return tasks
}

func printTimeNow() {
	t := time.Now()
	// 四舍五入到百分之一秒（两位小数）
	roundedNsec := (t.UnixNano() + 5e6) / 1e7 * 1e7 // 加5e6纳秒（0.5毫秒）实现四舍五入
	roundedTime := time.Unix(0, roundedNsec).In(t.Location())
	// 格式化为字符串，保留两位小数
	fmt.Print("\n当前时间为：")
	fmt.Print(roundedTime.Format("2006-01-02 15:04:05"))
	println()
}
