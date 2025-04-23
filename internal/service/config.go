package service

import (
	"fmt"
	"os"
	"time"

	"github.com/zhoucq/airflow-db-cleaner/internal/database"
	"github.com/zhoucq/airflow-db-cleaner/internal/models"
	"gopkg.in/yaml.v2"
)

// AppConfig 存储应用配置
type AppConfig struct {
	Database struct {
		Host            string        `yaml:"host"`
		Port            int           `yaml:"port"`
		User            string        `yaml:"user"`
		Password        string        `yaml:"password"`
		Name            string        `yaml:"name"`
		MaxIdleConns    int           `yaml:"max_idle_conns"`
		MaxOpenConns    int           `yaml:"max_open_conns"`
		ConnMaxLifetime time.Duration `yaml:"conn_max_lifetime"`
		Mock            bool          `yaml:"mock"`
	} `yaml:"database"`

	Cleaner struct {
		RetentionDays struct {
			DagRun       int `yaml:"dag_run"`
			TaskInstance int `yaml:"task_instance"`
			XCom         int `yaml:"xcom"`
			Log          int `yaml:"log"`
			Job          int `yaml:"job"`
		} `yaml:"retention_days"`
		BatchSize           int           `yaml:"batch_size"`
		SleepBetweenBatches time.Duration `yaml:"sleep_between_batches"`
		SleepSeconds        int           `yaml:"sleep_seconds"`
		DryRun              bool          `yaml:"dry_run"`
		Verbose             bool          `yaml:"verbose"`
	} `yaml:"cleaner"`

	Log struct {
		Level string `yaml:"level"`
		File  string `yaml:"file"`
	} `yaml:"log"`
}

// LoadConfig 从文件中加载配置
func LoadConfig(configPath string) (*AppConfig, error) {
	data, err := os.ReadFile(configPath)
	if err != nil {
		return nil, fmt.Errorf("读取配置文件失败: %w", err)
	}

	var config AppConfig
	if err := yaml.Unmarshal(data, &config); err != nil {
		return nil, fmt.Errorf("解析配置文件失败: %w", err)
	}

	// 设置默认值
	if config.Cleaner.BatchSize <= 0 {
		config.Cleaner.BatchSize = 1000
	}
	if config.Cleaner.SleepSeconds <= 0 {
		config.Cleaner.SleepSeconds = 5
	}

	return &config, nil
}

// GetDatabaseConfig 提取数据库配置
func (c *AppConfig) GetDatabaseConfig() database.Config {
	return database.Config{
		Host:            c.Database.Host,
		Port:            c.Database.Port,
		User:            c.Database.User,
		Password:        c.Database.Password,
		Name:            c.Database.Name,
		MaxIdleConns:    c.Database.MaxIdleConns,
		MaxOpenConns:    c.Database.MaxOpenConns,
		ConnMaxLifetime: c.Database.ConnMaxLifetime,
		Mock:            c.Database.Mock,
	}
}

// GetCleanerConfig 提取清理配置
func (c *AppConfig) GetCleanerConfig() models.Config {
	return models.Config{
		RetentionDays: map[string]int{
			"dag_run":       c.Cleaner.RetentionDays.DagRun,
			"task_instance": c.Cleaner.RetentionDays.TaskInstance,
			"xcom":          c.Cleaner.RetentionDays.XCom,
			"log":           c.Cleaner.RetentionDays.Log,
			"job":           c.Cleaner.RetentionDays.Job,
		},
		BatchSize:    c.Cleaner.BatchSize,
		DryRun:       c.Cleaner.DryRun,
		Verbose:      c.Cleaner.Verbose,
		SleepSeconds: c.Cleaner.SleepSeconds,
	}
}
