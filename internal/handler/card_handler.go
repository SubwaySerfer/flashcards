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

	var tagIDs []uuid.UUID
	for _, t := range card.Tags {
		tagIDs = append(tagIDs, t.ID)
	}

	if card.ID == uuid.Nil {
		card.ID = uuid.New()
	}

	if err := h.cardService.CreateCard(c.Request.Context(), &card, tagIDs); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, card)
}

func (h *CardHandler) GetCard(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid card ID"})
		return
	}

	card, err := h.cardService.GetCard(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, card)
}

func (h *CardHandler) UpdateCard(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid card ID"})
		return
	}

	var card domain.Card
	if err := c.ShouldBindJSON(&card); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	card.ID = id

	var tagIDs []uuid.UUID
	for _, t := range card.Tags {
		tagIDs = append(tagIDs, t.ID)
	}

	if err := h.cardService.UpdateCard(c.Request.Context(), &card, tagIDs); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, card)
}

func (h *CardHandler) DeleteCard(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid card ID"})
		return
	}

	if err := h.cardService.DeleteCard(c.Request.Context(), id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.Status(http.StatusNoContent)
}

func (h *CardHandler) ListCards(c *gin.Context) {
	cards, err := h.cardService.ListCards(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, cards)
}
