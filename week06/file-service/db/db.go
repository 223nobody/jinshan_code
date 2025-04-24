package db

import (
	"database/sql"
	"fmt"

	"github.com/jmoiron/sqlx"
	_ "modernc.org/sqlite"
)

type FileService struct {
	Id   int
	Name string
	Size int64
	Uuid string
	Type string
}

type Database struct {
	db *sqlx.DB
}

const (
	// SQL建表语句（使用TEXT类型存储序列化后的JSON字符串）
	createTableSQL = `CREATE TABLE file_service (
    id        INTEGER PRIMARY KEY AUTOINCREMENT,
    name      TEXT NOT NULL,        -- 原始文件名
    size      INTEGER CHECK(size > 0), 
    uuid      TEXT UNIQUE NOT NULL, -- 存储系统唯一ID
    type      TEXT NOT NULL        -- MIME类型
	);`
)

// 创建数据库
func InitDatabase(dsn string) (*Database, error) {
	db, err := sqlx.Connect("sqlite", dsn)
	if err != nil {
		return nil, fmt.Errorf("数据库连接失败: %w", err)
	}

	if _, err := db.Exec(createTableSQL); err != nil {
		return nil, fmt.Errorf("表创建失败: %w", err)
	}

	return &Database{db: db}, nil
}

func (d *Database) Close() error {
	return d.db.Close()
}

// Select_id 根据ID查询记录
func (d *Database) Select_id(id int) (*FileService, error) {
	query := `SELECT * FROM file_service WHERE id = ?`

	var file FileService
	if err := d.db.Get(&file, query, id); err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("记录不存在")
		}
		return nil, fmt.Errorf("查询失败: %w", err)
	}

	return &file, nil
}

// Select_type 根据Type查询记录
func (d *Database) Select_type(fileType string) (*FileService, error) {
	query := `SELECT * FROM file_service WHERE type = ?`
	var file FileService
	if err := d.db.Get(&file, query, fileType); err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("记录不存在")
		}
		return nil, fmt.Errorf("查询失败: %w", err)
	}

	return &file, nil
}

// Insert 插入记录
func (d *Database) Insert(name string, size int64, uuid string, fileType string) (int, error) {
	query := `INSERT INTO file_service (name, size, uuid, type) 
			  VALUES (:name, :size, :uuid, :type)`

	params := map[string]interface{}{
		"name": name,
		"size": size,
		"uuid": uuid,
		"type": fileType,
	}

	result, err := d.db.NamedExec(query, params)
	if err != nil {
		return 0, fmt.Errorf("创建失败: %w", err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("获取ID失败: %w", err)
	}

	return int(id), nil
}

// Update 更新记录
func (d *Database) Update(f *FileService) error {
	query := `UPDATE file_service SET 
              name = :name, 
              size = :size, 
              type = :type 
              WHERE id = :id`

	result, err := d.db.NamedExec(query, f)
	if err != nil {
		return fmt.Errorf("更新失败: %w", err)
	}

	if rows, _ := result.RowsAffected(); rows == 0 {
		return fmt.Errorf("无匹配记录")
	}

	return nil
}

// Delete 删除记录
func (d *Database) Delete(id int) error {
	query := `DELETE FROM file_service WHERE id = ?`

	result, err := d.db.Exec(query, id)
	if err != nil {
		return fmt.Errorf("删除失败: %w", err)
	}

	if rows, _ := result.RowsAffected(); rows == 0 {
		return fmt.Errorf("无匹配记录")
	}

	return nil
}
