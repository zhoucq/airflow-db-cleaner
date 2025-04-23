package database

import (
	"database/sql"
	"fmt"
	"log"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
)

// Config 数据库配置
type Config struct {
	Host            string
	Port            int
	User            string
	Password        string
	Name            string
	MaxIdleConns    int
	MaxOpenConns    int
	ConnMaxLifetime time.Duration
	Mock            bool // 模拟模式，不实际连接数据库
}

// DB 封装数据库连接
type DB struct {
	*sqlx.DB
	mock bool
}

// MockDB 是模拟数据库实现
type MockDB struct {
	*sqlx.DB
}

// New 创建数据库连接
func New(config Config) (*DB, error) {
	// 如果是模拟模式，返回模拟数据库实现
	if config.Mock {
		log.Printf("使用模拟模式，不实际连接数据库")
		return &DB{nil, true}, nil
	}

	// 构建DSN连接字符串
	var dsn string
	if config.Password == "" {
		// 空密码
		dsn = fmt.Sprintf("%s@tcp(%s:%d)/%s?parseTime=true&loc=Local",
			config.User, config.Host, config.Port, config.Name)
	} else {
		// 有密码
		dsn = fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?parseTime=true&loc=Local",
			config.User, config.Password, config.Host, config.Port, config.Name)
	}

	db, err := sqlx.Connect("mysql", dsn)
	if err != nil {
		return nil, fmt.Errorf("连接数据库失败: %w", err)
	}

	// 设置连接池参数
	db.SetMaxIdleConns(config.MaxIdleConns)
	db.SetMaxOpenConns(config.MaxOpenConns)
	db.SetConnMaxLifetime(config.ConnMaxLifetime)

	// 测试连接
	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("测试数据库连接失败: %w", err)
	}

	log.Printf("成功连接到数据库 %s:%d/%s", config.Host, config.Port, config.Name)
	return &DB{db, false}, nil
}

// Close 关闭数据库连接
func (db *DB) Close() error {
	if db.mock {
		return nil
	}
	return db.DB.Close()
}

// Get 获取单条记录
func (db *DB) Get(dest interface{}, query string, args ...interface{}) error {
	if db.mock {
		log.Printf("[模拟] 执行查询: %s, 参数: %v", query, args)

		// 模拟一些数据
		if intPtr, ok := dest.(*int); ok {
			*intPtr = 1000 // 模拟记录数
			return nil
		}

		return nil
	}
	return db.DB.Get(dest, query, args...)
}

// Select 获取多条记录
func (db *DB) Select(dest interface{}, query string, args ...interface{}) error {
	if db.mock {
		log.Printf("[模拟] 执行查询: %s, 参数: %v", query, args)

		// 如果查询主键，返回空结果
		if strSlice, ok := dest.(*[]string); ok {
			*strSlice = []string{"id"} // 模拟主键
			return nil
		}

		return nil
	}
	return db.DB.Select(dest, query, args...)
}

// Exec 执行SQL
func (db *DB) Exec(query string, args ...interface{}) (sql.Result, error) {
	if db.mock {
		log.Printf("[模拟] 执行SQL: %s, 参数: %v", query, args)
		return MockResult{1000}, nil
	}
	return db.DB.Exec(query, args...)
}

// Queryx 查询
func (db *DB) Queryx(query string, args ...interface{}) (*sqlx.Rows, error) {
	if db.mock {
		log.Printf("[模拟] 执行查询: %s, 参数: %v", query, args)
		return nil, fmt.Errorf("模拟模式不支持Queryx")
	}
	return db.DB.Queryx(query, args...)
}

// MockResult 是模拟结果
type MockResult struct {
	AffectedRows int64
}

// LastInsertId 实现Result接口
func (r MockResult) LastInsertId() (int64, error) {
	return 0, nil
}

// RowsAffected 实现Result接口
func (r MockResult) RowsAffected() (int64, error) {
	return r.AffectedRows, nil
}
