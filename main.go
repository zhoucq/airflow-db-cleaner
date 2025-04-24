package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/zhoucq/airflow-db-cleaner/internal/database"
	"github.com/zhoucq/airflow-db-cleaner/internal/service"
)

func main() {
	// Parse command line arguments
	configPath := flag.String("config", "config/config.yaml", "Configuration file path")
	flag.Parse()

	// Ensure the configuration file path is absolute
	absConfigPath, err := filepath.Abs(*configPath)
	if err != nil {
		log.Fatalf("Unable to get absolute path of configuration file: %v", err)
	}

	// Check if file exists
	if _, err := os.Stat(absConfigPath); os.IsNotExist(err) {
		log.Fatalf("Configuration file does not exist: %s", absConfigPath)
	}

	// Load configuration
	config, err := service.LoadConfig(absConfigPath)
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	// Set log according to configuration
	if config.Log.File != "" {
		logFile, err := os.OpenFile(config.Log.File, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
		if err != nil {
			log.Fatalf("Unable to open log file: %v", err)
		}
		defer logFile.Close()
		log.SetOutput(logFile)
	}

	// Connect to database
	db, err := database.New(config.GetDatabaseConfig())
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	// Create cleaner
	cleaner := service.NewCleaner(db, config.GetCleanerConfig())

	// Print run mode
	if config.Cleaner.DryRun {
		fmt.Println("=== Running in Dry Run mode ===")
		fmt.Println("No actual deletion operations will be executed, only showing the number of records to be deleted")
	} else {
		fmt.Println("=== Running in Execution mode ===")
		fmt.Println("Actual deletion operations will be executed, please ensure important data has been backed up")
	}

	// Print deletion method
	if config.Cleaner.UsePrimaryKeyDelete {
		fmt.Println("\n=== Using Primary Key-based deletion method ===")
		fmt.Println("This method can be faster for large tables but may involve more queries")
	} else {
		fmt.Println("\n=== Using Direct DELETE method ===")
		fmt.Println("This method is simpler but may be slower for large tables")
	}

	fmt.Println("\n=== Starting to clean expired data ===")

	// Execute cleaning
	if err := cleaner.CleanAll(); err != nil {
		log.Fatalf("Failed to clean data: %v", err)
	}

	fmt.Println("=== Data cleaning completed ===")
}
