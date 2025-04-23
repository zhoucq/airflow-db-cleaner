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
		{TableName: "dag_run", RetentionDays: c.config.RetentionDays["dag_run"], DateColumn: "execution_date"},
		{TableName: "task_instance", RetentionDays: c.config.RetentionDays["task_instance"], DateColumn: "start_date"},
		{TableName: "xcom", RetentionDays: c.config.RetentionDays["xcom"], DateColumn: "timestamp"},
		{TableName: "log", RetentionDays: c.config.RetentionDays["log"], DateColumn: "dttm"},
		{TableName: "job", RetentionDays: c.config.RetentionDays["job"], DateColumn: "end_date"},
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

	// 确保日期列存在
	var columnExists int
	checkColumnSQL := fmt.Sprintf(`
		SELECT COUNT(*) 
		FROM information_schema.columns 
		WHERE table_schema = DATABASE() 
		AND table_name = '%s' 
		AND column_name = '%s'
	`, table.TableName, table.DateColumn)

	if err := c.db.Get(&columnExists, checkColumnSQL); err != nil {
		return fmt.Errorf("检查列是否存在失败: %w", err)
	}

	if columnExists == 0 {
		log.Printf("警告: 表 %s 中不存在列 %s，跳过此表", table.TableName, table.DateColumn)
		return nil
	}

	// 先获取符合条件的记录数
	var count int
	countQuery := fmt.Sprintf("SELECT COUNT(*) FROM `%s` WHERE `%s` < ?", table.TableName, table.DateColumn)
	if err := c.db.Get(&count, countQuery, cutoffDate); err != nil {
		return fmt.Errorf("获取记录数失败: %w", err)
	}

	log.Printf("将清理表 %s 中的 %d 条记录", table.TableName, count)

	// 如果是演习模式，到此结束
	if c.config.DryRun {
		log.Printf("干运行模式: 不执行实际删除操作")
		return nil
	}

	// 如果没有记录需要清理，直接返回
	if count == 0 {
		log.Printf("表 %s 中没有过期记录需要清理", table.TableName)
		return nil
	}

	// 分批删除数据
	var deleted int
	batchSize := c.config.BatchSize
	sleepDuration := time.Duration(c.config.SleepSeconds) * time.Second

	// 使用简单的批量删除方式
	for deleted < count {
		// 计算本次批次要删除的记录数
		currentBatchSize := batchSize
		if count-deleted < batchSize {
			currentBatchSize = count - deleted
		}

		deleteSQL := fmt.Sprintf("DELETE FROM `%s` WHERE `%s` < ? LIMIT %d",
			table.TableName, table.DateColumn, currentBatchSize)

		result, err := c.db.Exec(deleteSQL, cutoffDate)
		if err != nil {
			return fmt.Errorf("删除记录失败: %w", err)
		}

		rowsAffected, _ := result.RowsAffected()
		deleted += int(rowsAffected)

		log.Printf("已删除表 %s 中的 %d/%d 条记录", table.TableName, deleted, count)

		// 如果没有删除完，休眠一段时间减轻数据库压力
		if deleted < count {
			log.Printf("休眠 %v 秒后继续删除...", sleepDuration.Seconds())
			time.Sleep(sleepDuration)
		}
	}

	log.Printf("成功清理表 %s 中的 %d 条记录", table.TableName, deleted)
	return nil
}
