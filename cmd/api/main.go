package main

import (
	"flashcards/internal/config"
	"flashcards/internal/domain"
	"flashcards/internal/handler"
	"flashcards/internal/repository"
	"flashcards/internal/service"
	"log"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// @title Flashcards API
// @version 1.0
// @description API Server for Flashcards application
// @host localhost:8080
// @BasePath /api/v1
func main() {
	// Initialize logger
	loggerConfig := zap.NewProductionConfig()
	loggerConfig.EncoderConfig.TimeKey = "timestamp"
	loggerConfig.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder

	logger, err := loggerConfig.Build()
	if err != nil {
		log.Fatalf("Cannot initialize logger: %v", err)
	}
	defer logger.Sync()

	cfg := config.NewConfig()
	logger.Info("Configuration loaded successfully",
		zap.String("server_port", cfg.Server.Port),
		zap.String("env", cfg.Environment))

	// Initialize DB
	dsn := cfg.Database.Url
	logger.Info("Attempting database connection",
		zap.String("database_url", dsn))

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		logger.Fatal("Failed to connect to database",
			zap.Error(err))
	}
	logger.Info("Database connected successfully")

	// Auto migrate the schema
	logger.Info("Starting database migration")
	err = db.AutoMigrate(&domain.Card{}, &domain.Tag{}, &domain.LearningProgress{})
	if err != nil {
		logger.Fatal("Failed to migrate database schema",
			zap.Error(err))
	}
	logger.Info("Database migration completed")

	// Initialize repositories
	logger.Debug("Initializing repositories")
	cardRepo := repository.NewCardRepository(db)
	tagRepo := repository.NewTagRepository(db)

	// Initialize services
	logger.Debug("Initializing services")
	cardService := service.NewCardService(cardRepo, tagRepo)
	tagService := service.NewTagService(tagRepo)

	// Initialize handlers
	logger.Debug("Initializing handlers")
	cardHandler := handler.NewCardHandler(cardService)
	tagHandler := handler.NewTagHandler(tagService)

	// Setup router
	logger.Info("Setting up router")
	r := gin.Default()

	// Add CORS middleware with custom configuration
	logger.Info("Configuring CORS",
		zap.Strings("allowed_origins", []string{"http://localhost:5173", "https://flashcards-service.netlify.app"}))

	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:5173", "https://flashcards-service.netlify.app"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	r.Use(func(c *gin.Context) {
		start := time.Now()
		path := c.Request.URL.Path

		c.Next()

		logger.Info("HTTP Request",
			zap.String("method", c.Request.Method),
			zap.String("path", path),
			zap.Int("status", c.Writer.Status()),
			zap.Duration("latency", time.Since(start)),
			zap.String("client_ip", c.ClientIP()),
			zap.String("user_agent", c.Request.UserAgent()),
		)
	})

	logger.Info("Setting up API routes")
	v1 := r.Group("/api/v1")
	{
		cards := v1.Group("/cards")
		{
			cards.POST("/", cardHandler.CreateCard)
			cards.GET("/:id", cardHandler.GetCard)
			cards.PUT("/:id", cardHandler.UpdateCard)
			cards.DELETE("/:id", cardHandler.DeleteCard)
			cards.GET("/", cardHandler.ListCards)
			cards.GET("/random", cardHandler.GetRandomCard)
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
	serverAddr := ":" + cfg.Server.Port
	logger.Info("Starting server",
		zap.String("address", serverAddr))

	if err := r.Run(serverAddr); err != nil {
		logger.Fatal("Server failed to start",
			zap.Error(err))
	}
}
