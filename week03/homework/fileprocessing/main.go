package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"os"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

// 数据结构体定义
type Word struct {
	Word         string        `json:"word"`
	Translations []Translation `json:"translations"`
	Phrases      []Phrase      `json:"phrases,omitempty"` //phrases可以为空
}

type Translation struct {
	Translation string `json:"translation"`
	Type        string `json:"type"`
}

type Phrase struct {
	Phrase      string `json:"phrase"`
	Translation string `json:"translation"`
}

const (
	DbName1        = "orderdCET4.db"
	DbName2        = "orderdCET6.db"
	filename1      = "3-CET4-顺序.json"
	filename2      = "4-CET6-顺序.json"
	createTableSQL = `CREATE TABLE IF NOT EXISTS words (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		word TEXT NOT NULL UNIQUE,
		translations TEXT NOT NULL,
		phrases TEXT
	);`
)

func main() {
	_ = os.Remove(DbName1)
	_ = os.Remove(DbName2)
	start := time.Now()

	//新建数据库db1,db2
	db1, err := initNewDatabase(DbName1)
	if err != nil {
		panic(err)
	}
	defer db1.Close()

	db2, err := initNewDatabase(DbName2)
	if err != nil {
		panic(err)
	}
	defer db2.Close()

	// 加载并插入数据
	words1, err := loadJSONData(filename1)
	if err != nil {
		panic(err)
	}

	if err := insertData(db1, words1); err != nil {
		panic(err)
	}

	fmt.Printf("成功新建数据库 %s 并插入 %d 条数据！\n", DbName1, len(words1))
	fmt.Printf("%s 插入用时: %.2f 秒\n", filename1, time.Since(start).Seconds())

	words2, err := loadJSONData(filename2)
	if err != nil {
		panic(err)
	}

	if err := insertData(db2, words2); err != nil {
		panic(err)
	}

	fmt.Printf("\n成功新建数据库 %s 并插入 %d 条数据！\n", DbName2, len(words2))
	fmt.Printf("%s 插入用时: %.2f 秒\n", filename2, time.Since(start).Seconds())
}

// 初始化新数据库
func initNewDatabase(DbName string) (*sql.DB, error) {
	// 创建新数据库文件
	dbName := DbName
	db, err := sql.Open("sqlite3", dbName)
	if err != nil {
		return nil, fmt.Errorf("创建数据库失败: %v", err)
	}

	// 执行建表语句
	if _, err := db.Exec(createTableSQL); err != nil {
		return nil, fmt.Errorf("创建表失败: %v", err)
	}

	return db, nil
}

// 读取json数据
func loadJSONData(filename string) ([]Word, error) {
	Filename := filename
	file, err := os.ReadFile(Filename)
	if err != nil {
		return nil, fmt.Errorf("读取JSON失败: %v", err)
	}

	var words []Word
	if err := json.Unmarshal(file, &words); err != nil {
		return nil, fmt.Errorf("解析JSON失败: %v", err)
	}

	return words, nil
}

// 导入json数据
func insertData(db *sql.DB, words []Word) error {
	tx, err := db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()
	//忽略可能的存在相同word的记录
	stmt, err := tx.Prepare(`INSERT OR IGNORE INTO words 
		(word, translations, phrases) 
		VALUES (?, ?, ?)`)
	if err != nil {
		return err
	}
	defer stmt.Close()

	//检验json各部分是否成功完成序列化
	for _, word := range words {
		transJSON, err := json.Marshal(word.Translations)
		if err != nil {
			return fmt.Errorf("序列化翻译失败: %v", err)
		}

		var phrasesJSON []byte
		if len(word.Phrases) > 0 {
			phrasesJSON, err = json.Marshal(word.Phrases)
			if err != nil {
				return fmt.Errorf("序列化短语失败: %v", err)
			}
		}

		if _, err := stmt.Exec(word.Word, string(transJSON), phrasesJSON); err != nil {
			return fmt.Errorf("插入失败: %v", err)
		}
	}

	return tx.Commit()
}
