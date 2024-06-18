package domain

import (
	"github.com/aerosystems/checkmail-service/internal/presenters/http/handlers"
)

type Handler struct {
	*handlers.BaseHandler
	domainUsecase handlers.DomainUsecase
}

func NewHandler(baseHandler *handlers.BaseHandler, domainUsecase handlers.DomainUsecase) *Handler {
	return &Handler{
		BaseHandler:   baseHandler,
		domainUsecase: domainUsecase,
	}
}
