package review

import "github.com/aerosystems/checkmail-service/internal/presenters/http/handlers"

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
