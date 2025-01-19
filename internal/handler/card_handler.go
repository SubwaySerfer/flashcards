package handler

import (
	"flashcards/internal/domain"
	"flashcards/internal/service"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type CardHandler struct {
	cardService service.CardService
}

func NewCardHandler(cardService service.CardService) *CardHandler {
	return &CardHandler{cardService: cardService}
}

// CreateCard godoc
// @Summary Create a new card
// @Description Create a new flashcard with title and description
// @Tags cards
// @Accept json
// @Produce json
// @Param card body domain.Card true "Card object"
// @Success 201 {object} domain.Card
// @Router /cards [post]
func (h *CardHandler) CreateCard(c *gin.Context) {
	var card domain.Card
	if err := c.ShouldBindJSON(&card); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if card.ID == uuid.Nil {
		card.ID = uuid.New()
	}

	if err := h.cardService.CreateCard(c.Request.Context(), &card); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, card)
}
