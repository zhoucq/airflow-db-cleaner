package database

import (
	"database/sql"
	"fmt"
	"log"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
)

// Config database configuration
type Config struct {
	Host            string
	Port            int
	User            string
	Password        string
	Name            string
	MaxIdleConns    int
	MaxOpenConns    int
	ConnMaxLifetime time.Duration
	Mock            bool // Mock mode, does not actually connect to the database
}

// DB encapsulates database connection
type DB struct {
	*sqlx.DB
	mock bool
}

// MockDB is a mock database implementation
type MockDB struct {
	*sqlx.DB
}

// New creates a database connection
func New(config Config) (*DB, error) {
	// If in mock mode, return a mock database implementation
	if config.Mock {
		log.Printf("Using mock mode, not actually connecting to the database")
		return &DB{nil, true}, nil
	}

	// Build DSN connection string
	var dsn string
	if config.Password == "" {
		// Empty password
		dsn = fmt.Sprintf("%s@tcp(%s:%d)/%s?parseTime=true&loc=Local",
			config.User, config.Host, config.Port, config.Name)
	} else {
		// With password
		dsn = fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?parseTime=true&loc=Local",
			config.User, config.Password, config.Host, config.Port, config.Name)
	}

	db, err := sqlx.Connect("mysql", dsn)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	// Set connection pool parameters
	db.SetMaxIdleConns(config.MaxIdleConns)
	db.SetMaxOpenConns(config.MaxOpenConns)
	db.SetConnMaxLifetime(config.ConnMaxLifetime)

	// Test connection
	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to test database connection: %w", err)
	}

	log.Printf("Successfully connected to database %s:%d/%s", config.Host, config.Port, config.Name)
	return &DB{db, false}, nil
}

// Close closes the database connection
func (db *DB) Close() error {
	if db.mock {
		return nil
	}
	return db.DB.Close()
}

// Get retrieves a single record
func (db *DB) Get(dest interface{}, query string, args ...interface{}) error {
	if db.mock {
		log.Printf("[Mock] Execute query: %s, args: %v", query, args)

		// Mock some data
		if intPtr, ok := dest.(*int); ok {
			*intPtr = 1000 // Mock record count
			return nil
		}

		return nil
	}
	return db.DB.Get(dest, query, args...)
}

// Select retrieves multiple records
func (db *DB) Select(dest interface{}, query string, args ...interface{}) error {
	if db.mock {
		log.Printf("[Mock] Execute query: %s, args: %v", query, args)

		// If querying primary key, return empty result
		if strSlice, ok := dest.(*[]string); ok {
			*strSlice = []string{"id"} // Mock primary key
			return nil
		}

		return nil
	}
	return db.DB.Select(dest, query, args...)
}

// Exec executes SQL
func (db *DB) Exec(query string, args ...interface{}) (sql.Result, error) {
	if db.mock {
		log.Printf("[Mock] Execute SQL: %s, args: %v", query, args)
		return MockResult{1000}, nil
	}
	return db.DB.Exec(query, args...)
}

// Queryx queries
func (db *DB) Queryx(query string, args ...interface{}) (*sqlx.Rows, error) {
	if db.mock {
		log.Printf("[Mock] Execute query: %s, args: %v", query, args)
		return nil, fmt.Errorf("mock mode does not support Queryx")
	}
	return db.DB.Queryx(query, args...)
}

// MockResult is a mock result
type MockResult struct {
	AffectedRows int64
}

// LastInsertId implements the Result interface
func (r MockResult) LastInsertId() (int64, error) {
	return 0, nil
}

// RowsAffected implements the Result interface
func (r MockResult) RowsAffected() (int64, error) {
	return r.AffectedRows, nil
}
