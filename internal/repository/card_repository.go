package repository

import (
	"context"
	"flashcards/internal/domain"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type CardRepository interface {
	Create(ctx context.Context, card *domain.Card) error
	GetByID(ctx context.Context, id uuid.UUID) (*domain.Card, error)
	Update(ctx context.Context, card *domain.Card) error
	Delete(ctx context.Context, id uuid.UUID) error
	List(ctx context.Context) ([]domain.Card, error)
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

func (r *cardRepository) GetByID(ctx context.Context, id uuid.UUID) (*domain.Card, error) {
	var card domain.Card
	err := r.db.WithContext(ctx).First(&card, "id = ?", id).Error
	return &card, err
}

func (r *cardRepository) Update(ctx context.Context, card *domain.Card) error {
	return r.db.WithContext(ctx).Save(card).Error
}

func (r *cardRepository) Delete(ctx context.Context, id uuid.UUID) error {
	return r.db.WithContext(ctx).Delete(&domain.Card{}, "id = ?", id).Error
}

func (r *cardRepository) List(ctx context.Context) ([]domain.Card, error) {
	var cards []domain.Card
	err := r.db.WithContext(ctx).Find(&cards).Error
	return cards, err
}
