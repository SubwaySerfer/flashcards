package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	Server   ServerConfig
	Database DatabaseConfig
}

type ServerConfig struct {
	Port string `default:"8080"`
}

type DatabaseConfig struct {
	Host     string `default:"localhost"`
	Port     string `default:"5432"`
	User     string `default:"postgres"`
	Password string
	DBName   string `default:"flashcards"`
	SSLMode  string `default:"disable"`
}

func NewConfig() *Config {
	err := godotenv.Load()
	if err != nil {
		log.Printf("Warning: No .env file found or it could not be loaded: %v", err)
	}

	return &Config{
		Server: ServerConfig{
			Port: "8081",
		},
		Database: DatabaseConfig{
			Host:     os.Getenv("DB_HOST"),
			Port:     os.Getenv("DB_PORT"),
			User:     os.Getenv("DB_USER"),
			Password: os.Getenv("DB_PASSWORD"),
			DBName:   os.Getenv("DB_NAME"),
			SSLMode:  os.Getenv("DB_SSLMODE"),
		},
	}
}
