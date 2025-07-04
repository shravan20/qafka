package main

import (
	"context"
	"log"

	"github.com/shravan20/qafka/internal/api"
	"github.com/shravan20/qafka/internal/config"
	"github.com/shravan20/qafka/internal/database"
	"github.com/shravan20/qafka/internal/services"

	"github.com/go-fuego/fuego"
	"github.com/joho/godotenv"
)

// @title Qafka API
// @version 1.0
// @description Universal Queue Management Platform API
// @termsOfService http://swagger.io/terms/

// @contact.name API Support
// @contact.url http://github.com/shravan20/qafka
// @contact.email support@qafka.io

// @license.name MIT
// @license.url https://opensource.org/licenses/MIT

// @host localhost:8080
// @BasePath /api/v1
func main() {
	// Load environment variables
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using system environment variables")
	}

	// Load configuration
	cfg := config.Load()

	// Initialize database
	db, err := database.Connect(cfg.DatabaseURL)
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}
	defer db.Close()

	// Run migrations
	if err := database.RunMigrations(context.Background(), db); err != nil {
		log.Fatal("Failed to run migrations:", err)
	}

	// Initialize services
	queueService := services.NewQueueService(db)
	monitoringService := services.NewMonitoringService()

	// Create Fuego app
	app := fuego.NewServer(
		fuego.WithPort(cfg.APIPort),
		fuego.WithCORS(fuego.CORSConfig{
			AllowOrigins: cfg.CORSOrigins,
			AllowMethods: []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
			AllowHeaders: []string{"*"},
		}),
	)

	// Setup routes
	api.SetupRoutes(app, queueService, monitoringService)

	// Setup Swagger documentation
	api.SetupSwagger(app)

	// Start server
	log.Printf("Starting Qafka API server on port %s", cfg.APIPort)
	log.Printf("Swagger docs available at http://localhost:%s/swagger/", cfg.APIPort)

	if err := app.Run(); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}
