package db

import (
	"auth-service/configs"
	"fmt"
	"log"
	"os"

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

	// Execute SQL statements
	_, err = sqlDB.Exec(sqlStatements)
	if err != nil {
		return fmt.Errorf("failed to execute SQL file: %v", err)
	}

	return nil
}