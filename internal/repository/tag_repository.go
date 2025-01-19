package repository

import (
	"context"
	"flashcards/internal/domain"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type TagRepository interface {
	Create(ctx context.Context, tag *domain.Tag) error
	GetByID(ctx context.Context, id uuid.UUID) (*domain.Tag, error)
	Update(ctx context.Context, tag *domain.Tag) error
	List(ctx context.Context) ([]domain.Tag, error)
}

type tagRepository struct {
	db *gorm.DB
}

func NewTagRepository(db *gorm.DB) TagRepository {
	return &tagRepository{db: db}
}

func (r *tagRepository) Create(ctx context.Context, tag *domain.Tag) error {
	if tag.ID == uuid.Nil {
		tag.ID = uuid.New()
	}
	return r.db.WithContext(ctx).Create(tag).Error
}

func (r *tagRepository) GetByID(ctx context.Context, id uuid.UUID) (*domain.Tag, error) {
	var tag domain.Tag
	err := r.db.WithContext(ctx).First(&tag, "id = ?", id).Error
	return &tag, err
}

func (r *tagRepository) Update(ctx context.Context, tag *domain.Tag) error {
	return r.db.WithContext(ctx).Save(tag).Error
}

func (r *tagRepository) List(ctx context.Context) ([]domain.Tag, error) {
	var tags []domain.Tag
	err := r.db.WithContext(ctx).Find(&tags).Error
	return tags, err
}
