package main

import (
	"flashcards/internal/config"
	"flashcards/internal/domain"
	"flashcards/internal/handler"
	"flashcards/internal/repository"
	"flashcards/internal/service"
	"fmt"
	"log"

	"github.com/gin-gonic/gin"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// @title Flashcards API
// @version 1.0
// @description API Server for Flashcards application
// @host localhost:8080
// @BasePath /api/v1
func main() {
	cfg := config.NewConfig()

	// Initialize DB
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=%s",
		cfg.Database.Host,
		cfg.Database.User,
		cfg.Database.Password,
		cfg.Database.DBName,
		cfg.Database.Port,
		cfg.Database.SSLMode,
	)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	// Auto migrate the schema
	db.AutoMigrate(&domain.Card{}, &domain.Tag{}, &domain.LearningProgress{})

	// Initialize repositories
	cardRepo := repository.NewCardRepository(db)

	// Initialize services
	cardService := service.NewCardService(cardRepo)

	// Initialize handlers
	cardHandler := handler.NewCardHandler(cardService)

	// Setup router
	r := gin.Default()

	v1 := r.Group("/api/v1")
	{
		cards := v1.Group("/cards")
		{
			cards.POST("/", cardHandler.CreateCard)
			// Add other routes here
		}
	}

	// Start server
	if err := r.Run(":" + cfg.Server.Port); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}
