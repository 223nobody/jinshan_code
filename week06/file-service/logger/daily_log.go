package logger

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sync"
	"time"
)

type DailyLogger struct {
	mu         sync.Mutex
	file       *os.File
	currentDay string
	logDir     string
}

func NewDailyLogger(logDir string) *DailyLogger {
	_ = os.MkdirAll(logDir, 0755)
	return &DailyLogger{
		logDir: logDir,
	}
}

func (dl *DailyLogger) getLogFile() (*os.File, error) {
	now := time.Now()
	today := now.Format("2006_01_02")

	if dl.currentDay != today || dl.file == nil {
		if dl.file != nil {
			_ = dl.file.Close()
		}

		filename := fmt.Sprintf("%s.json", today)
		path := filepath.Join(dl.logDir, filename)
		f, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE, 0644)
		if err != nil {
			return nil, err
		}

		dl.file = f
		dl.currentDay = today
	}

	return dl.file, nil
}

func (dl *DailyLogger) Log(data map[string]interface{}) error {
	dl.mu.Lock()
	defer dl.mu.Unlock()

	file, err := dl.getLogFile()
	if err != nil {
		return err
	}

	// 准备读取现有内容
	if _, err = file.Seek(0, 0); err != nil {
		return err
	}

	var entries []map[string]interface{}
	decoder := json.NewDecoder(file)
	if err = decoder.Decode(&entries); err != nil {
		// 处理空文件或无效内容
		if err == io.EOF || err == io.ErrUnexpectedEOF {
			entries = make([]map[string]interface{}, 0)
		} else {
			return err
		}
	}

	// 添加新条目
	entries = append(entries, data)

	// 清空并重写文件
	if err = file.Truncate(0); err != nil {
		return err
	}
	if _, err = file.Seek(0, 0); err != nil {
		return err
	}

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ") // 添加格式化缩进
	if err = encoder.Encode(entries); err != nil {
		return err
	}

	return nil
}
