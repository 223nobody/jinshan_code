package storage

import (
	"encoding/json"
	"fmt"
	"sort"

	"github.com/jmoiron/sqlx"
	_ "modernc.org/sqlite"

	"Server/config"
)

const createTableSQL = `
CREATE TABLE IF NOT EXISTS questions (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    title TEXT NOT NULL,
    type INTEGER NOT NULL,
	language TEXT NOT NULL,
    answers TEXT COMMENT NOT NULL,
	rights TEXT COMMENT NOT NULL
);
`

// Database 包装器结构体
type Database struct {
	db *sqlx.DB
}

// InitDB 返回自定义 Database 类型
func InitDB(dsn string) (*Database, error) {
	db, err := sqlx.Connect("sqlite", dsn)
	if err != nil {
		return nil, fmt.Errorf("数据库连接失败: %w", err)
	}

	// 执行建表语句
	if _, err := db.Exec(createTableSQL); err != nil {
		return nil, fmt.Errorf("初始化表失败: %w", err)
	}

	return &Database{db: db}, nil
}

// 补充数据库操作方法
func (d *Database) Close() error {
	return d.db.Close()
}

// 在 Database 结构体中添加事务方法
func (d *Database) Beginx() (*sqlx.Tx, error) {
	return d.db.Beginx()
}

// Select 查询多条记录
func (d *Database) Select(dest interface{}, query string, args ...interface{}) error {
	return d.db.Select(dest, query, args...)
}

// Get 查询单条记录
func (d *Database) Get(dest interface{}, query string, args ...interface{}) error {

	return d.db.Get(dest, query, args...)

}

// Exec 删除多条数据
func (d *Database) Exec(query string, args ...interface{}) error {
	_, err := d.db.Exec(query, args...)
	return err
}

// 添加数据库操作
func (d *Database) CreateQuestion(q *config.QuestionRequest1) (int, error) {
	const query = `
        INSERT INTO questions 
        (title,type ,language, answers, rights)
        VALUES ($1, $2, $3, $4, $5)
        RETURNING id`

	answersJSON, _ := json.Marshal(q.Answers)
	rightsJSON, _ := json.Marshal(q.Rights)

	var id int
	err := d.db.QueryRow(
		query,
		q.Title,
		q.Type,
		q.Language,
		string(answersJSON),
		string(rightsJSON),
	).Scan(&id)

	return id, err
}

// UpdateQuestion 更新题目
func (db *Database) UpdateQuestion(req *config.QuestionRequest1) (int64, error) {
	// 1. 数据格式转换（匹配图片中的选项结构）
	optionsJSON, err := json.Marshal(req.Answers)
	if err != nil {
		return 0, fmt.Errorf("选项序列化失败: %v", err)
	}

	// 2. 答案排序（匹配图片中的A/B/D顺序）
	sort.Strings(req.Rights)
	answersJSON, err := json.Marshal(req.Rights)
	if err != nil {
		return 0, fmt.Errorf("答案序列化失败: %v", err)
	}

	// 3. 执行SQL更新（包含图片中的所有字段）
	result, err := db.db.Exec(`
        UPDATE questions SET
            title = ?,
            type = ?,
            language = ?,
            answers = ?,
            rights = ?
        WHERE id = ?`,
		req.Title,
		req.Type,
		req.Language,
		string(optionsJSON),
		string(answersJSON),
		req.Id,
	)

	if err != nil {
		return 0, fmt.Errorf("SQL执行失败: %v", err)
	}

	// 4. 返回影响行数（匹配图片中的编辑场景）
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return 0, fmt.Errorf("获取影响行数失败: %v", err)
	}
	return rowsAffected, nil
}
