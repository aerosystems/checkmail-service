package review

import (
	"github.com/aerosystems/checkmail-service/internal/models"
	"github.com/aerosystems/checkmail-service/internal/presenters/http/handlers"
	"time"
)

type Handler struct {
	*handlers.BaseHandler
	reviewUsecase handlers.ReviewUsecase
}

func NewHandler(
	baseHandler *handlers.BaseHandler,
	reviewUsecase handlers.ReviewUsecase,
) *Handler {
	return &Handler{
		BaseHandler:   baseHandler,
		reviewUsecase: reviewUsecase,
	}
}

type Review struct {
	Name      string    `json:"name" example:"gmail.com"`
	Type      string    `json:"type" example:"whitelist"`
	CreatedAt time.Time `json:"createdAt" example:"2021-01-01T00:00:00Z"`
	UpdatedAt time.Time `json:"updatedAt" example:"2021-01-01T00:00:00Z"`
}

func ModelToReview(review models.Review) Review {
	return Review{
		Name:      review.Name,
		Type:      review.Type,
		CreatedAt: review.CreatedAt,
		UpdatedAt: review.UpdatedAt,
	}
}
