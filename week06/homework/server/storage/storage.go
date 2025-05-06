package storage

import (
	"Server/config"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sync"
	"time"
)

type Storage interface {
	Save(log config.AILog) error
}

type JSONStorage struct {
	basePath string
	mu       sync.Mutex
}

func NewJSONStorage() *JSONStorage {
	return &JSONStorage{
		basePath: "log",
	}
}

func (s *JSONStorage) Save(log config.AILog) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	// 创建目录
	if err := os.MkdirAll(s.basePath, 0755); err != nil {
		return fmt.Errorf("创建目录失败: %w", err)
	}

	// 生成文件名
	filename := time.Now().Format("2006_01_02") + ".json"
	filepath := filepath.Join(s.basePath, filename)

	// 读取现有数据
	var logs []config.AILog
	if fileExists(filepath) {
		data, err := os.ReadFile(filepath)
		if err != nil {
			return fmt.Errorf("读取文件失败: %w", err)
		}

		// 如果文件非空则解析
		if len(data) > 0 {
			if err := json.Unmarshal(data, &logs); err != nil {
				return fmt.Errorf("解析现有数据失败: %w", err)
			}
		}
	}

	// 追加新记录
	logs = append(logs, log)

	// 写入更新后的数据
	file, err := os.Create(filepath) // 使用Create覆盖写入
	if err != nil {
		return fmt.Errorf("创建文件失败: %w", err)
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ") // 保持美观格式
	if err := encoder.Encode(logs); err != nil {
		return fmt.Errorf("数据编码失败: %w", err)
	}

	return nil
}

// 检查文件是否存在
func fileExists(path string) bool {
	_, err := os.Stat(path)
	return !os.IsNotExist(err)
}
