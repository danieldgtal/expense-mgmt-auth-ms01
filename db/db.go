package db

import (
	"auth-service/configs"
	"fmt"
	"log"
	"os"
	"strings"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// DB - Global database connection
var DB *gorm.DB


// ConnectDatabase initializes a connection to the PostgreSQL database
func ConnectDatabase() error {
	config := configs.LoadConfig() // Load the configuration

	// Build the connection string
	dbURI := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
		config.Database.Host,
		config.Database.Port,
		config.Database.User,
		config.Database.Password,
		config.Database.Name,
		config.Database.SSLMode)

	// Open a connection to the database
	var err error
	DB, err = gorm.Open(postgres.Open(dbURI), &gorm.Config{})
	if err != nil {
		return fmt.Errorf("failed to initialize database: %w", err)
	}

	log.Println("Database connected successfully")
	return nil // Return nil error if successful
}

// ExecuteSQLSchema reads a SQL file and executes its content
func ExecuteSQLSchema(filePath string) error {
	// Open a raw database connection using sql.DB
	sqlDB, err := DB.DB()
	if err != nil {
		return fmt.Errorf("failed to get raw database connection: %v", err)
	}

	// Read SQL schema file
	sqlBytes, err := os.ReadFile(filePath)
	if err != nil {
		return fmt.Errorf("failed to read SQL file: %v", err)
	}
	sqlStatements := string(sqlBytes)

	// Split SQL statements by semicolon
	statements := strings.Split(sqlStatements, ";")

	// Begin transaction
	tx, err := sqlDB.Begin()
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %v", err)
	}
	defer tx.Rollback() // Ensure rollback if not committed

	// Execute each statement
	for _, stmt := range statements {
		stmt = strings.TrimSpace(stmt) // Clean up whitespace
		if stmt == "" {
			continue // Skip empty statements
		}

		_, err := tx.Exec(stmt)
		if err != nil {
			return fmt.Errorf("failed to execute SQL statement: %v - Error: %v", stmt, err)
		}
		log.Printf("Executed SQL statement: %s\n", stmt)
	}

	// Commit the transaction
	if err := tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit transaction: %v", err)
	}

	return nil
}

// GetDBInstance returns the DB instance
func GetDBInstance() *gorm.DB {
	return DB
}