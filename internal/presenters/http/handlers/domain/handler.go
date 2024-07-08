package domain

import (
	"github.com/aerosystems/checkmail-service/internal/models"
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

type Domain struct {
	Name     string `json:"name" example:"gmail.com"`
	Type     string `json:"type" example:"whitelist"`
	Coverage string `json:"coverage" example:"equals"`
}

func ModelToDomain(model *models.Domain) Domain {
	return Domain{
		Name:     model.Name,
		Type:     model.Type.String(),
		Coverage: model.Coverage.String(),
	}
}
