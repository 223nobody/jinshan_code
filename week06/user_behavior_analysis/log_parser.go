package main

import (
	"bufio"
	"encoding/csv"
	"fmt"
	"os"
	"strings"
	"time"
)

// LogRecord 对应日志文件中的单条记录
type LogRecord struct {
	Timestamp time.Time // 时间戳应为第一个字段
	UserID    string
	Action    string
	Detail    string // 新增第四个字段
}

// 定义标准时间格式
const (
	timeFormat = "2006-01-02 15:04:05" // 匹配日志中的时间格式
	fieldCount = 4                     // 每行应有4个字段
)

// ParseLogFile 解析日志文件
func ParseLogFile(filename string) ([]LogRecord, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, fmt.Errorf("打开文件失败: %w", err)
	}
	defer file.Close()

	var records []LogRecord
	scanner := bufio.NewScanner(file)
	lineNumber := 0 // 从0开始计数更准确

	for scanner.Scan() {
		lineNumber++
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue // 跳过空行
		}

		record, err := parseLogLine(line)
		if err != nil {
			return nil, fmt.Errorf("第%d行解析失败: %w", lineNumber, err)
		}
		records = append(records, record)
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("文件读取错误: %w", err)
	}

	return records, nil
}

// 更新后的解析函数
func parseLogLine(line string) (LogRecord, error) {
	// 使用带逗号分割的csv解析方式
	reader := csv.NewReader(strings.NewReader(line))
	fields, err := reader.Read()
	if err != nil {
		return LogRecord{}, fmt.Errorf("CSV解析失败: %w", err)
	}

	// 验证字段数量
	if len(fields) != fieldCount {
		return LogRecord{}, fmt.Errorf("需要%d个字段，实际收到%d个", fieldCount, len(fields))
	}

	// 解析时间戳（第一个字段）
	timestamp, err := time.Parse(timeFormat, strings.TrimSpace(fields[0]))
	if err != nil {
		return LogRecord{}, fmt.Errorf("时间解析失败: %w", err)
	}

	return LogRecord{
		Timestamp: timestamp,
		UserID:    strings.TrimSpace(fields[1]),
		Action:    strings.TrimSpace(fields[2]),
		Detail:    strings.TrimSpace(fields[3]),
	}, nil
}
