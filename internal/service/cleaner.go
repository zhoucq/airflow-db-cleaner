package service

import (
	"fmt"
	"log"
	"time"

	"github.com/zhoucq/airflow-db-cleaner/internal/database"
	"github.com/zhoucq/airflow-db-cleaner/internal/models"
)

// Cleaner responsible for cleaning expired data
type Cleaner struct {
	db     *database.DB
	config models.Config
}

// NewCleaner creates a new cleaner
func NewCleaner(db *database.DB, config models.Config) *Cleaner {
	return &Cleaner{
		db:     db,
		config: config,
	}
}

// CleanAll cleans all configured tables
func (c *Cleaner) CleanAll() error {
	// Prepare table configurations for cleaning
	tables := []models.TableConfig{
		{TableName: "dag_run", RetentionDays: c.config.RetentionDays["dag_run"], DateColumn: "execution_date"},
		{TableName: "task_instance", RetentionDays: c.config.RetentionDays["task_instance"], DateColumn: "start_date"},
		{TableName: "xcom", RetentionDays: c.config.RetentionDays["xcom"], DateColumn: "timestamp"},
		{TableName: "log", RetentionDays: c.config.RetentionDays["log"], DateColumn: "dttm"},
		{TableName: "job", RetentionDays: c.config.RetentionDays["job"], DateColumn: "end_date"},
	}

	// Iterate and clean each table
	for _, table := range tables {
		if err := c.cleanTable(table); err != nil {
			return fmt.Errorf("failed to clean table %s: %w", table.TableName, err)
		}
	}

	return nil
}

// cleanTable cleans expired data from the specified table
func (c *Cleaner) cleanTable(table models.TableConfig) error {
	// Calculate cutoff date
	cutoffDate := time.Now().AddDate(0, 0, -table.RetentionDays)
	log.Printf("Preparing to clean table %s with data earlier than %s", table.TableName, cutoffDate.Format("2006-01-02"))

	// Ensure date column exists
	var columnExists int
	checkColumnSQL := fmt.Sprintf(`
		SELECT COUNT(*) 
		FROM information_schema.columns 
		WHERE table_schema = DATABASE() 
		AND table_name = '%s' 
		AND column_name = '%s'
	`, table.TableName, table.DateColumn)

	if err := c.db.Get(&columnExists, checkColumnSQL); err != nil {
		return fmt.Errorf("failed to check if column exists: %w", err)
	}

	if columnExists == 0 {
		log.Printf("Warning: Column %s does not exist in table %s, skipping this table", table.DateColumn, table.TableName)
		return nil
	}

	// First get the number of records that match the condition
	var count int
	countQuery := fmt.Sprintf("SELECT COUNT(*) FROM `%s` WHERE `%s` < ?", table.TableName, table.DateColumn)
	if err := c.db.Get(&count, countQuery, cutoffDate); err != nil {
		return fmt.Errorf("failed to get record count: %w", err)
	}

	log.Printf("Will clean %d records from table %s", count, table.TableName)

	// If in dry run mode, stop here
	if c.config.DryRun {
		log.Printf("Dry run mode: No actual deletion operations will be performed")
		return nil
	}

	// If there are no records to clean, return directly
	if count == 0 {
		log.Printf("No expired records need to be cleaned in table %s", table.TableName)
		return nil
	}

	// Delete data in batches
	var deleted int
	batchSize := c.config.BatchSize
	sleepDuration := time.Duration(c.config.SleepSeconds) * time.Second

	// Use simple batch deletion method
	for deleted < count {
		// Calculate the number of records to delete in this batch
		currentBatchSize := batchSize
		if count-deleted < batchSize {
			currentBatchSize = count - deleted
		}

		deleteSQL := fmt.Sprintf("DELETE FROM `%s` WHERE `%s` < ? LIMIT %d",
			table.TableName, table.DateColumn, currentBatchSize)

		result, err := c.db.Exec(deleteSQL, cutoffDate)
		if err != nil {
			return fmt.Errorf("failed to delete records: %w", err)
		}

		rowsAffected, _ := result.RowsAffected()
		deleted += int(rowsAffected)

		log.Printf("Deleted %d/%d records from table %s", deleted, count, table.TableName)

		// If not finished deleting, sleep to reduce database pressure
		if deleted < count {
			log.Printf("Sleeping for %v seconds before continuing deletion...", sleepDuration.Seconds())
			time.Sleep(sleepDuration)
		}
	}

	log.Printf("Successfully cleaned %d records from table %s", deleted, table.TableName)
	return nil
}
