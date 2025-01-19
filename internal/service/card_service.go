package service

import (
	"context"
	"database/sql"
	"flashcards/internal/domain"
	"flashcards/internal/repository"
	"fmt"
	"os"

	"github.com/google/uuid"
	_ "github.com/lib/pq"
)

type CardService interface {
	GetCard(ctx context.Context, id uuid.UUID) (*domain.Card, error)
	UpdateCard(ctx context.Context, card *domain.Card) error
	DeleteCard(ctx context.Context, id uuid.UUID) error
	ListCards(ctx context.Context) ([]domain.Card, error)
	CreateCard(ctx context.Context, card *domain.Card, tagIDs []uuid.UUID) error
	CreateDatabase(ctx context.Context) error
}

type cardService struct {
	repo    repository.CardRepository
	tagRepo repository.TagRepository
}

func NewCardService(repo repository.CardRepository, tagRepo repository.TagRepository) CardService {
	return &cardService{repo: repo, tagRepo: tagRepo}
}

// func (s *cardService) CreateCard(ctx context.Context, card *domain.Card) error {
// 	return s.repo.Create(ctx, card)
// }

func (s *cardService) CreateCard(ctx context.Context, card *domain.Card, tagIDs []uuid.UUID) error {
	var tags []domain.Tag
	if err := s.repo.GetTagsByIds(ctx, tagIDs, &tags); err != nil {
		return err
	}
	card.Tags = tags
	return s.repo.Create(ctx, card)
}

func (s *cardService) GetCard(ctx context.Context, id uuid.UUID) (*domain.Card, error) {
	return s.repo.GetByID(ctx, id)
}

func (s *cardService) UpdateCard(ctx context.Context, card *domain.Card) error {
	return s.repo.Update(ctx, card)
}

func (s *cardService) DeleteCard(ctx context.Context, id uuid.UUID) error {
	return s.repo.Delete(ctx, id)
}

func (s *cardService) ListCards(ctx context.Context) ([]domain.Card, error) {
	return s.repo.ListCards(ctx)
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
