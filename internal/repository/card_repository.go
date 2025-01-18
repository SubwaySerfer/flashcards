package repository

import (
	"context"
	"flashcards/internal/domain"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type CardRepository interface {
	Create(ctx context.Context, card *domain.Card) error
	//GetByID(ctx context.Context, id uuid.UUID) (*domain.Card, error)
	//Update(ctx context.Context, card *domain.Card) error
	//Delete(ctx context.Context, id uuid.UUID) error
	//List(ctx context.Context, filters map[string]interface{}) ([]domain.Card, error)
	//FindByTags(ctx context.Context, tags []string) ([]domain.Card, error)
}

type cardRepository struct {
	db *gorm.DB
}

func NewCardRepository(db *gorm.DB) CardRepository {
	return &cardRepository{db: db}
}

func (r *cardRepository) Create(ctx context.Context, card *domain.Card) error {
	if card.ID == uuid.Nil {
		card.ID = uuid.New()
	}
	return r.db.WithContext(ctx).Create(card).Error
}
