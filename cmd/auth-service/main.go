package main

import (
	"log"
	"os"

	"github.com/Debt-Solvers/BE-auth-service/db"
	"github.com/Debt-Solvers/BE-auth-service/internal/middleware"
	"github.com/Debt-Solvers/BE-auth-service/internal/routes"

	"github.com/gin-gonic/gin"
)


func main() {
	// Initialize database connection
	if err := db.ConnectDatabase(); err != nil {
		log.Fatalf("Database connection error: %v", err)
	}

	// Execute Schema
	if err := db.ExecuteSQLSchema("./db/schema.sql"); err != nil {
		log.Fatalf("Error executing schema: %v", err)
	}

	// Initialize Gin engine
	server := gin.Default()
	// Register the logging middleware
	server.Use(middleware.Logger()) 

	// Register routes
	routes.RegisterRoutes(server) // Public routes
	routes.AddHealthCheckRoute(server)

	// Check for environment variable port
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	// Start the server
	if err := server.Run(":" + port); err != nil {
		log.Fatalf("Failed to start the server: %v", err)
	}
}