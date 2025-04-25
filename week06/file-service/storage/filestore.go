package storage

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"time"
)

type FileStore struct {
	basePath string
}

func NewFileStore(basePath string) *FileStore {
	return &FileStore{basePath}
}

// 智能路径生成（防止单目录文件过多）
func (fs *FileStore) generatePath(uuid string) string {
	dateDir := time.Now().Format("2006-01-02")
	subDir := filepath.Join(fs.basePath, dateDir, uuid[:2])
	_ = os.MkdirAll(subDir, 0755)
	return filepath.Join(subDir, uuid)
}

func (fs *FileStore) Save(uuid string, src io.Reader) error {
	dstPath := fs.generatePath(uuid)
	file, err := os.Create(dstPath)
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = io.Copy(file, src)
	return err
}

func (fs *FileStore) Get(uuid string) (io.ReadCloser, error) {
	return os.Open(fs.findFile(uuid))
}

// 递归查找文件（考虑日期目录结构）
func (fs *FileStore) findFile(uuid string) string {
	// 限制递归深度（最多查找最近365天的目录）
	maxDepth := 365
	daysAgo := 0

	// 构造可能的UUID前缀子目录（取前2位）
	prefix := uuid[:2]
	if len(prefix) < 2 {
		prefix = "00" // 处理短UUID情况
	}

	// 从当天开始向前回溯查找
	for i := 0; i < maxDepth; i++ {
		// 计算目标日期目录（格式：YYYY-MM-DD）
		targetDate := time.Now().AddDate(0, 0, -daysAgo)
		dateDir := targetDate.Format("2006-01-02")

		// 构建完整文件路径
		targetPath := filepath.Join(
			fs.basePath,
			dateDir,
			prefix,
			uuid,
		)

		// 检查文件是否存在
		if _, err := os.Stat(targetPath); err == nil {
			return targetPath
		}

		// 当天未找到则回溯前一天
		daysAgo++
	}

	// 全量遍历模式（当快速查找失败时）
	return fs.deepSearch(uuid, prefix)
}

// 深度搜索实现（遍历所有日期目录）
func (fs *FileStore) deepSearch(uuid, prefix string) string {
	filepath.Walk(fs.basePath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return nil // 跳过错误目录
		}

		// 匹配目标文件模式：*/<prefix>/<uuid>
		matched, _ := filepath.Match(
			filepath.Join("*", prefix, uuid),
			path,
		)

		if matched {
			// 找到文件路径后通过panic提前返回
			panic(path)
		}
		return nil
	})

	// 捕获panic获取结果
	defer func() {
		if r := recover(); r != nil {
			if path, ok := r.(string); ok {
				panic(path) // 继续抛出结果
			}
		}
	}()

	return "" // 未找到返回空
}

// Delete 删除文件
func (fs *FileStore) Delete(uuid string) error {
	filePath := fs.findFile(uuid)
	if filePath == "" {
		return fmt.Errorf("文件不存在")
	}

	if err := os.Remove(filePath); err != nil {
		return fmt.Errorf("删除失败: %w", err)
	}

	// 清理空目录（可选）
	dirPath := filepath.Dir(filePath)
	if isEmpty, _ := isDirEmpty(dirPath); isEmpty {
		_ = os.Remove(dirPath)
	}

	return nil
}

// 检查目录是否为空
func isDirEmpty(path string) (bool, error) {
	f, err := os.Open(path)
	if err != nil {
		return false, err
	}
	defer f.Close()

	_, err = f.Readdirnames(1)
	return err == io.EOF, nil
}
