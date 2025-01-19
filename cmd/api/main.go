package main

import (
	"flashcards/internal/config"
	"flashcards/internal/domain"
	"flashcards/internal/handler"
	"flashcards/internal/repository"
	"flashcards/internal/service"
	"fmt"
	"log"

	"github.com/gin-contrib/cors"
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
	tagRepo := repository.NewTagRepository(db)

	// Initialize services
	cardService := service.NewCardService(cardRepo, tagRepo)
	tagService := service.NewTagService(tagRepo)

	// Initialize handlers
	cardHandler := handler.NewCardHandler(cardService)
	tagHandler := handler.NewTagHandler(tagService)

	// Setup router
	r := gin.Default()

	// Add CORS middleware
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:5173"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
	}))
	r.Use(gin.Logger())
	r.Use(func(c *gin.Context) {
		log.Printf("Request: %s %s", c.Request.Method, c.Request.URL)
		c.Next()
	})

	v1 := r.Group("/api/v1")
	{
		cards := v1.Group("/cards")
		{
			cards.POST("/", cardHandler.CreateCard)
			cards.GET("/:id", cardHandler.GetCard)
			cards.PUT("/:id", cardHandler.UpdateCard)
			cards.DELETE("/:id", cardHandler.DeleteCard)
			cards.GET("/", cardHandler.ListCards)
		}
		tags := v1.Group("/tags")
		{
			tags.POST("/", tagHandler.CreateTag)
			tags.GET("/:id", tagHandler.GetTag)
			tags.PUT("/:id", tagHandler.UpdateTag)
			tags.GET("/", tagHandler.ListTags)
		}
	}

	// Start server
	if err := r.Run(":" + cfg.Server.Port); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}
