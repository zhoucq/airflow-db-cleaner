package service

import (
	"fmt"
	"log"
	"strings"
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
		{TableName: "dag_run", RetentionDays: c.config.RetentionDays["dag_run"], DateColumn: "execution_date", PrimaryKey: "id"},
		{TableName: "task_instance", RetentionDays: c.config.RetentionDays["task_instance"], DateColumn: "start_date", PrimaryKey: "dag_id,task_id,run_id,map_index"},
		{TableName: "xcom", RetentionDays: c.config.RetentionDays["xcom"], DateColumn: "timestamp", PrimaryKey: "dag_id,task_id,run_id,map_index,key"},
		{TableName: "log", RetentionDays: c.config.RetentionDays["log"], DateColumn: "dttm", PrimaryKey: "id"},
		{TableName: "job", RetentionDays: c.config.RetentionDays["job"], DateColumn: "end_date", PrimaryKey: "id"},
	}

	// Iterate and clean each table
	for _, table := range tables {
		var err error
		if c.config.UsePrimaryKeyDelete {
			err = c.cleanTableByPK(table)
		} else {
			err = c.cleanTable(table)
		}

		if err != nil {
			return fmt.Errorf("failed to clean table %s: %w", table.TableName, err)
		}
	}

	return nil
}

// cleanTable cleans expired data from the specified table using the original method
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
	// Convert float64 seconds to time.Duration (nanoseconds)
	sleepDuration := time.Duration(c.config.SleepSeconds * float64(time.Second))

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
			log.Printf("Sleeping for %.3f seconds before continuing deletion...", c.config.SleepSeconds)
			time.Sleep(sleepDuration)
		}
	}

	log.Printf("Successfully cleaned %d records from table %s", deleted, table.TableName)
	return nil
}

// cleanTableByPK cleans expired data from the specified table using primary key-based deletion
func (c *Cleaner) cleanTableByPK(table models.TableConfig) error {
	// Calculate cutoff date
	cutoffDate := time.Now().AddDate(0, 0, -table.RetentionDays)
	log.Printf("Preparing to clean table %s with data earlier than %s (using PK-based method)",
		table.TableName, cutoffDate.Format("2006-01-02"))

	// Validate required fields
	if table.PrimaryKey == "" {
		return fmt.Errorf("primary key not specified for table %s", table.TableName)
	}

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
	sleepDuration := time.Duration(c.config.SleepSeconds * float64(time.Second))
	pk := strings.Split(table.PrimaryKey, ",")

	for deleted < count {
		// Calculate actual batch size for this iteration
		currentBatchSize := batchSize
		if count-deleted < batchSize {
			currentBatchSize = count - deleted
		}

		startTime := time.Now()

		// First query: get primary keys of records to delete
		// If multi-column primary key, select them all
		pkSelectSQL := fmt.Sprintf("SELECT %s FROM `%s` WHERE `%s` < ? ORDER BY %s LIMIT %d",
			table.PrimaryKey, table.TableName, table.DateColumn, pk[0], currentBatchSize)

		// For composite keys, we'll read the rows into a map of string->interface{}
		rows, err := c.db.Queryx(pkSelectSQL, cutoffDate)
		if err != nil {
			return fmt.Errorf("failed to query primary keys: %w", err)
		}

		// Prepare for batch deletion
		var batchDeleted int

		// Handle different primary key scenarios
		if len(pk) == 1 {
			// Single column primary key
			var ids []interface{}
			for rows.Next() {
				var id interface{}
				if err := rows.Scan(&id); err != nil {
					rows.Close()
					return fmt.Errorf("failed to scan primary key: %w", err)
				}
				ids = append(ids, id)
			}
			rows.Close()

			if len(ids) == 0 {
				break // No more records to delete
			}

			// Build placeholders for IN clause
			placeholders := strings.Repeat("?,", len(ids))
			placeholders = placeholders[:len(placeholders)-1] // Remove trailing comma

			// Delete using primary key
			deleteSQL := fmt.Sprintf("DELETE FROM `%s` WHERE `%s` IN (%s)",
				table.TableName, pk[0], placeholders)

			result, err := c.db.Exec(deleteSQL, ids...)
			if err != nil {
				return fmt.Errorf("failed to delete records: %w", err)
			}

			rowsAffected, _ := result.RowsAffected()
			batchDeleted = int(rowsAffected)
		} else {
			// Composite primary key - handle differently
			// We'll do one deletion per row to avoid complex WHERE clauses

			type CompositeKey struct {
				Values []interface{}
			}

			var keys []CompositeKey
			for rows.Next() {
				// Scan into a map
				rowMap := make(map[string]interface{})
				if err := rows.MapScan(rowMap); err != nil {
					rows.Close()
					return fmt.Errorf("failed to scan composite key: %w", err)
				}

				// Extract values in the right order
				var keyValues []interface{}
				for _, col := range pk {
					keyValues = append(keyValues, rowMap[col])
				}

				keys = append(keys, CompositeKey{Values: keyValues})
			}
			rows.Close()

			if len(keys) == 0 {
				break // No more records to delete
			}

			// Build WHERE clause for composite key
			// For example: (col1 = ? AND col2 = ? AND col3 = ?) OR (col1 = ? AND col2 = ? AND col3 = ?) ...
			var whereClauseParts []string
			var allParams []interface{}

			for i, key := range keys {
				var conditions []string
				for j, col := range pk {
					conditions = append(conditions, fmt.Sprintf("`%s` = ?", col))
					allParams = append(allParams, key.Values[j])
				}
				whereClauseParts = append(whereClauseParts, "("+strings.Join(conditions, " AND ")+")")

				// For very large batches, limit the size of a single DELETE query
				if len(whereClauseParts) >= 100 || i == len(keys)-1 {
					whereClause := strings.Join(whereClauseParts, " OR ")
					deleteSQL := fmt.Sprintf("DELETE FROM `%s` WHERE %s", table.TableName, whereClause)

					result, err := c.db.Exec(deleteSQL, allParams...)
					if err != nil {
						return fmt.Errorf("failed to delete records with composite key: %w", err)
					}

					rowsAffected, _ := result.RowsAffected()
					batchDeleted += int(rowsAffected)

					// Reset for next batch
					whereClauseParts = nil
					allParams = nil
				}
			}
		}

		deleted += batchDeleted

		// Calculate execution time for this batch
		batchDuration := time.Since(startTime)
		log.Printf("Deleted %d/%d records from table %s (batch time: %.2fs)",
			deleted, count, table.TableName, batchDuration.Seconds())

		// If not finished deleting, sleep to reduce database pressure
		if deleted < count {
			log.Printf("Sleeping for %.3f seconds before continuing deletion...", c.config.SleepSeconds)
			time.Sleep(sleepDuration)
		}
	}

	log.Printf("Successfully cleaned %d records from table %s", deleted, table.TableName)
	return nil
}
