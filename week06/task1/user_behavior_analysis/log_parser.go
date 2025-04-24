package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"time"
)

// LogRecord 对应日志文件中的单条记录
type LogRecord struct {
	UserID    string
	Action    string
	Timestamp time.Time
}

// ParseLogFile 解析日志文件
func ParseLogFile(filename string) ([]LogRecord, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, fmt.Errorf("打开文件失败: %w", err)
	}
	defer file.Close()

	var records []LogRecord
	scanner := bufio.NewScanner(file)
	lineNumber := 1

	for scanner.Scan() {
		line := scanner.Text()
		record, err := parseLogLine(line)
		if err != nil {
			return nil, fmt.Errorf("第%d行解析失败: %w", lineNumber, err)
		}
		records = append(records, record)
		lineNumber++
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("文件读取错误: %w", err)
	}

	return records, nil
}

// 示例日志格式：user001,login,2023-08-20T09:15:00Z
func parseLogLine(line string) (LogRecord, error) {
	parts := strings.Split(line, ",")
	if len(parts) != 3 {
		return LogRecord{}, fmt.Errorf("无效的日志格式")
	}

	timestamp, err := time.Parse(time.RFC3339, parts[2])
	if err != nil {
		return LogRecord{}, fmt.Errorf("时间解析失败: %w", err)
	}

	return LogRecord{
		UserID:    parts[0],
		Action:    parts[1],
		Timestamp: timestamp,
	}, nil
}
