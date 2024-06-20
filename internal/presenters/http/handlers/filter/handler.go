package filter

import (
	"github.com/aerosystems/checkmail-service/internal/presenters/http/handlers"
)

type Handler struct {
	*handlers.BaseHandler
	filterUsecase handlers.FilterUsecase
}

func NewHandler(
	baseHandler *handlers.BaseHandler,
	filterUsecase handlers.FilterUsecase,
) *Handler {
	return &Handler{
		BaseHandler:   baseHandler,
		filterUsecase: filterUsecase,
	}
}
