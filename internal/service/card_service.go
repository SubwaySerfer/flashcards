package service

import (
	"context"
	"database/sql"
	"flashcards/internal/domain"
	"flashcards/internal/repository"
	"fmt"
	"os"

	_ "github.com/google/uuid"
	_ "github.com/lib/pq"
)

type CardService interface {
	CreateCard(ctx context.Context, card *domain.Card) error
	//GetCard(ctx context.Context, id uuid.UUID) (*domain.Card, error)
	//UpdateCard(ctx context.Context, card *domain.Card) error
	//DeleteCard(ctx context.Context, id uuid.UUID) error
	//ListCards(ctx context.Context, filters map[string]interface{}) ([]domain.Card, error)
	//FindCardsByTags(ctx context.Context, tags []string) ([]domain.Card, error)
	CreateDatabase(ctx context.Context) error
}

type cardService struct {
	repo repository.CardRepository
}

func NewCardService(repo repository.CardRepository) CardService {
	return &cardService{repo: repo}
}

func (s *cardService) CreateCard(ctx context.Context, card *domain.Card) error {
	return s.repo.Create(ctx, card)
}

func (s *cardService) CreateDatabase(ctx context.Context) error {
	println("Creating database", os.Getenv("DB_USER"), os.Getenv("DB_PASSWORD"), os.Getenv("DB_HOST"), os.Getenv("DB_PORT"), os.Getenv("DB_SSLMODE"))
	connStr := fmt.Sprintf("user=%s password=%s host=%s port=%s sslmode=%s",
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_HOST"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_SSLMODE"),
	)
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return fmt.Errorf("failed to connect to PostgreSQL: %w", err)
	}
	defer db.Close()

	_, err = db.ExecContext(ctx, "CREATE DATABASE flashcards")
	if err != nil {
		return fmt.Errorf("failed to create database: %w", err)
	}

	return nil
}
