package domain

import (
	"time"

	"github.com/google/uuid"
)

type Card struct {
	ID          uuid.UUID `gorm:"type:uuid;primary_key" json:"id"`
	Title       string    `gorm:"size:255;not null" json:"title" validate:"required,min=1,max=255"`
	Description string    `gorm:"type:text" json:"description"`
	CreatedAt   time.Time `gorm:"not null" json:"created_at"`
	UpdatedAt   time.Time `gorm:"not null" json:"updated_at"`
	Tags        []Tag     `gorm:"many2many:card_tags;" json:"tags"`
}

type Tag struct {
	ID    uuid.UUID `gorm:"type:uuid;primary_key" json:"id"`
	Name  string    `gorm:"size:100;not null;unique" json:"name" validate:"required,min=1,max=100"`
	Cards []Card    `gorm:"many2many:card_tags;" json:"-"`
}

type LearningProgress struct {
	ID              uuid.UUID `gorm:"type:uuid;primary_key" json:"id"`
	CardID          uuid.UUID `gorm:"type:uuid;not null" json:"card_id"`
	Card            Card      `gorm:"foreignKey:CardID" json:"card"`
	LastReviewedAt  time.Time `gorm:"not null" json:"last_reviewed_at"`
	NextReviewAt    time.Time `gorm:"not null" json:"next_review_at"`
	DifficultyLevel int       `gorm:"not null" json:"difficulty_level" validate:"required,min=1,max=5"`
	TimesReviewed   int       `gorm:"not null;default:0" json:"times_reviewed"`
	SuccessRate     float64   `gorm:"not null;default:0" json:"success_rate"`
}
