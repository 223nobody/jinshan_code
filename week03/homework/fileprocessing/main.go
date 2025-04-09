package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"os"
	"time"

	_ "github.com/mattn/go-sqlite3" // 匿名导入SQLite驱动
)

type Word struct {
	Word         string        `json:"word"`              // 单词
	Translations []Translation `json:"translations"`      // 翻译切片（非空）
	Phrases      []Phrase      `json:"phrases,omitempty"` // 短语切片（可为空）
}

// Translation 单词的翻译内容及其词性
type Translation struct {
	Translation string `json:"translation"` // 翻译内容
	Type        string `json:"type"`        // 词性
}

// Phrase 包含短语及其翻译
type Phrase struct {
	Phrase      string `json:"phrase"`      // 短语
	Translation string `json:"translation"` // 短语翻译
}

const (
	DbName1   = "orderdCET4.db"  // CET4数据库文件名
	DbName2   = "orderdCET6.db"  // CET6数据库文件名
	filename1 = "3-CET4-顺序.json" // CET4源数据文件
	filename2 = "4-CET6-顺序.json" // CET6源数据文件

	// SQL建表语句（使用TEXT类型存储序列化后的JSON字符串）
	createTableSQL = `CREATE TABLE IF NOT EXISTS words (
		id INTEGER PRIMARY KEY AUTOINCREMENT,  -- 自增主键
		word TEXT NOT NULL UNIQUE,             -- 单词(唯一约束)
		translations TEXT NOT NULL,            -- 翻译JSON数组
		phrases TEXT                           -- 短语JSON数组(可为空)
	);`
)

func main() {
	// 删除旧数据库文件（如果存在）
	_ = os.Remove(DbName1)
	_ = os.Remove(DbName2)

	start := time.Now() // 记录程序启动时间

	// 初始化CET4数据库
	db1, err := initNewDatabase(DbName1)
	if err != nil {
		panic(fmt.Sprintf("初始化CET4数据库失败: %v", err))
	}
	defer db1.Close() // 主函数退出前关闭连接

	// 初始化CET6数据库
	db2, err := initNewDatabase(DbName2)
	if err != nil {
		panic(fmt.Sprintf("初始化CET6数据库失败: %v", err))
	}
	defer db2.Close()

	// 处理CET4数据
	words1, err := loadJSONData(filename1)
	if err != nil {
		panic(fmt.Sprintf("加载CET4数据失败: %v", err))
	}

	if err := insertData(db1, words1); err != nil {
		panic(fmt.Sprintf("插入CET4数据失败: %v", err))
	}

	fmt.Printf("[成功] 创建 %s | 插入 %d 条数据 | 耗时: %.2f秒\n",
		DbName1, len(words1), time.Since(start).Seconds())

	// 处理CET6数据
	words2, err := loadJSONData(filename2)
	if err != nil {
		panic(fmt.Sprintf("加载CET6数据失败: %v", err))
	}

	if err := insertData(db2, words2); err != nil {
		panic(fmt.Sprintf("插入CET6数据失败: %v", err))
	}

	fmt.Printf("\n[成功] 创建 %s | 插入 %d 条数据 | 总耗时: %.2f秒\n",
		DbName2, len(words2), time.Since(start).Seconds())
}

// initNewDatabase 初始化新数据库连接并创建表结构
// 返回值:
//
//	*sql.DB: 数据库连接对象
//	error: 错误信息
func initNewDatabase(DbName string) (*sql.DB, error) {
	// 创建数据库连接（如果文件不存在会直接创建）
	db, err := sql.Open("sqlite3", DbName)
	if err != nil {
		return nil, fmt.Errorf("连接数据库失败: %v", err)
	}

	// 执行建表语句
	if _, err := db.Exec(createTableSQL); err != nil {
		return nil, fmt.Errorf("执行建表语句失败: %v", err)
	}

	return db, nil
}

// loadJSONData 从指定JSON文件加载单词数据
// 返回值:
//
//	[]Word: 解析后的单词结构体切片
//	error: 错误信息
func loadJSONData(filename string) ([]Word, error) {
	// 读取整个文件内容
	fileContent, err := os.ReadFile(filename)
	if err != nil {
		return nil, fmt.Errorf("读取文件失败: %v", err)
	}

	var words []Word
	// 反序列化JSON到结构体切片
	if err := json.Unmarshal(fileContent, &words); err != nil {
		return nil, fmt.Errorf("解析JSON失败: %v", err)
	}

	return words, nil
}

// insertData 将单词数据批量插入数据库
// 传入:
//
//	db: 数据库连接对象
//	words: 待插入的单词切片
//
// 返回值:
//
//	error: 错误信息
func insertData(db *sql.DB, words []Word) error {
	// 开启事务（提升批量插入性能）
	tx, err := db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback() // 确保出错时回滚

	// 准备预处理语句（使用INSERT OR IGNORE避免重复插入）
	stmt, err := tx.Prepare(`INSERT OR IGNORE INTO words 
		(word, translations, phrases) 
		VALUES (?, ?, ?)`)
	if err != nil {
		return err
	}
	defer stmt.Close()

	// 遍历所有单词
	for _, word := range words {
		// 序列化翻译数据
		transJSON, err := json.Marshal(word.Translations)
		if err != nil {
			return fmt.Errorf("序列化翻译失败[%s]: %v", word.Word, err)
		}

		// 处理可选短语字段
		var phrasesJSON []byte
		if len(word.Phrases) > 0 {
			if phrasesJSON, err = json.Marshal(word.Phrases); err != nil {
				return fmt.Errorf("序列化短语失败[%s]: %v", word.Word, err)
			}
		}

		// 执行插入操作
		if _, err := stmt.Exec(word.Word, string(transJSON), phrasesJSON); err != nil {
			return fmt.Errorf("插入失败[%s]: %v", word.Word, err)
		}
	}
	// 提交事务
	return tx.Commit()
}
