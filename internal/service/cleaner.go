package service

import (
	"fmt"
	"log"
	"time"

	"github.com/zhoucq/airflow-db-cleaner/internal/database"
	"github.com/zhoucq/airflow-db-cleaner/internal/models"
)

// Cleaner 负责清理过期数据
type Cleaner struct {
	db     *database.DB
	config models.Config
}

// NewCleaner 创建新的清理器
func NewCleaner(db *database.DB, config models.Config) *Cleaner {
	return &Cleaner{
		db:     db,
		config: config,
	}
}

// CleanAll 清理所有配置的表
func (c *Cleaner) CleanAll() error {
	// 准备清理表的配置
	tables := []models.TableConfig{
		{TableName: "dag_run", RetentionDays: c.config.RetentionDays["dag_run"]},
		{TableName: "task_instance", RetentionDays: c.config.RetentionDays["task_instance"]},
		{TableName: "xcom", RetentionDays: c.config.RetentionDays["xcom"]},
		{TableName: "log", RetentionDays: c.config.RetentionDays["log"]},
		{TableName: "job", RetentionDays: c.config.RetentionDays["job"]},
	}

	// 遍历并清理每个表
	for _, table := range tables {
		if err := c.cleanTable(table); err != nil {
			return fmt.Errorf("清理表 %s 失败: %w", table.TableName, err)
		}
	}

	return nil
}

// cleanTable 清理指定表的过期数据
func (c *Cleaner) cleanTable(table models.TableConfig) error {
	// 计算截止日期
	cutoffDate := time.Now().AddDate(0, 0, -table.RetentionDays)
	log.Printf("准备清理表 %s 中早于 %s 的数据", table.TableName, cutoffDate.Format("2006-01-02"))

	// 构建查询条件
	var dateColumn string
	switch table.TableName {
	case "dag_run":
		dateColumn = "execution_date"
	case "task_instance":
		dateColumn = "execution_date"
	case "xcom":
		dateColumn = "execution_date"
	case "log":
		dateColumn = "created_at"
	case "job":
		dateColumn = "end_date"
	default:
		return fmt.Errorf("未知表: %s", table.TableName)
	}

	// 先获取符合条件的记录数
	var count int
	countQuery := fmt.Sprintf("SELECT COUNT(*) FROM %s WHERE %s < ?", table.TableName, dateColumn)
	if err := c.db.Get(&count, countQuery, cutoffDate); err != nil {
		return fmt.Errorf("获取记录数失败: %w", err)
	}

	log.Printf("将清理表 %s 中的 %d 条记录", table.TableName, count)

	// 如果是演习模式，到此结束
	if c.config.DryRun {
		log.Printf("干运行模式: 不执行实际删除操作")
		return nil
	}

	// 分批删除数据
	var deleted int
	batchSize := c.config.BatchSize
	for deleted < count {
		deleteQuery := fmt.Sprintf("DELETE FROM %s WHERE %s < ? LIMIT %d",
			table.TableName, dateColumn, batchSize)

		result, err := c.db.Exec(deleteQuery, cutoffDate)
		if err != nil {
			return fmt.Errorf("删除记录失败: %w", err)
		}

		rowsAffected, _ := result.RowsAffected()
		deleted += int(rowsAffected)

		log.Printf("已删除表 %s 中的 %d/%d 条记录", table.TableName, deleted, count)

		// 休眠一段时间，减轻数据库压力
		if deleted < count && c.config.BatchSize > 0 {
			sleepTime := time.Second
			time.Sleep(sleepTime)
		}
	}

	log.Printf("成功清理表 %s 中的 %d 条记录", table.TableName, deleted)
	return nil
}
