package HTTPServer

import (
	"github.com/aerosystems/checkmail-service/internal/models"
)

type DomainHandler struct {
	*BaseHandler
	domainUsecase DomainUsecase
}

func NewDomainHandler(baseHandler *BaseHandler, domainUsecase DomainUsecase) *DomainHandler {
	return &DomainHandler{
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
