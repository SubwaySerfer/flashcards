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
	ListCards(ctx context.Context) ([]domain.Card, error)
	GetTagsByIds(ctx context.Context, ids []uuid.UUID, tags *[]domain.Tag) error
	//FindByTags(ctx context.Context, tags []string) ([]domain.Card, error)
	GetRandomCard(ctx context.Context) (*domain.Card, error)
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
	err := r.db.WithContext(ctx).
		Preload("Tags").
		First(&card, "id = ?", id).Error
	return &card, err
}

func (r *cardRepository) Update(ctx context.Context, card *domain.Card) error {
	// Start a transaction
	tags := card.Tags
	tx := r.db.WithContext(ctx).Begin()

	// Delete old tags associations
	if err := tx.Model(&card).Association("Tags").Clear(); err != nil {
		tx.Rollback()
		return err
	}
	// Save the card with new tags
	card.Tags = tags
	if err := tx.Save(card).Error; err != nil {
		tx.Rollback()
		return err
	}
	// Reload the card with tags
	if err := tx.Preload("Tags").First(&card, "id = ?", card.ID).Error; err != nil {
		tx.Rollback()
		return err
	}

	// Commit the transaction
	return tx.Commit().Error
}

func (r *cardRepository) Delete(ctx context.Context, id uuid.UUID) error {
	return r.db.WithContext(ctx).Delete(&domain.Card{}, "id = ?", id).Error
}

func (r *cardRepository) ListCards(ctx context.Context) ([]domain.Card, error) {
	var cards []domain.Card
	if err := r.db.WithContext(ctx).
		Preload("Tags").
		Find(&cards).Error; err != nil {
		return nil, err
	}
	return cards, nil
}

func (r *cardRepository) GetTagsByIds(ctx context.Context, ids []uuid.UUID, tags *[]domain.Tag) error {
	return r.db.WithContext(ctx).Where("id IN ?", ids).Find(tags).Error
}

func (r *cardRepository) GetRandomCard(ctx context.Context) (*domain.Card, error) {
	var card domain.Card
	err := r.db.WithContext(ctx).
		Preload("Tags").
		Order("RANDOM()").
		First(&card).Error
	return &card, err
}
