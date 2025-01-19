package service

import (
	"context"
	"flashcards/internal/domain"
	"flashcards/internal/repository"

	"github.com/google/uuid"
)

type TagService interface {
	CreateTag(ctx context.Context, tag *domain.Tag) error
	GetTag(ctx context.Context, id uuid.UUID) (*domain.Tag, error)
	UpdateTag(ctx context.Context, tag *domain.Tag) error
	ListTags(ctx context.Context) ([]domain.Tag, error)
}

type tagService struct {
	repo repository.TagRepository
}

func NewTagService(repo repository.TagRepository) TagService {
	return &tagService{repo: repo}
}

func (s *tagService) CreateTag(ctx context.Context, tag *domain.Tag) error {
	return s.repo.Create(ctx, tag)
}

func (s *tagService) GetTag(ctx context.Context, id uuid.UUID) (*domain.Tag, error) {
	return s.repo.GetByID(ctx, id)
}

func (s *tagService) UpdateTag(ctx context.Context, tag *domain.Tag) error {
	return s.repo.Update(ctx, tag)
}

func (s *tagService) ListTags(ctx context.Context) ([]domain.Tag, error) {
	return s.repo.List(ctx)
}
