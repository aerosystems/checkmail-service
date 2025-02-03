package HTTPServer

import (
	"github.com/aerosystems/checkmail-service/internal/models"
	"time"
)

type Handler struct {
	accessUsecase  AccessUsecase
	inspectUsecase InspectUsecase
	domainUsecase  ManageUsecase
}

func NewHandler(
	accessUsecase AccessUsecase,
	inspectUsecase InspectUsecase,
	domainUsecase ManageUsecase,
) *Handler {
	return &Handler{
		accessUsecase:  accessUsecase,
		inspectUsecase: inspectUsecase,
		domainUsecase:  domainUsecase,
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
		Coverage: model.Match.String(),
	}
}

type Filter struct {
	Name      string    `json:"name" example:"gmail.com"`
	Type      string    `json:"type" example:"whitelist"`
	Match     string    `json:"coverage" example:"equals"`
	CreatedAt time.Time `json:"createdAt" example:"2024-01-01T00:00:00Z"`
	UpdatedAt time.Time `json:"updatedAt" example:"2024-01-01T00:00:00Z"`
}

func ModelToFilter(filter models.Filter) Filter {
	return Filter{
		Name:      filter.Name,
		Type:      filter.Type.String(),
		Match:     filter.Match.String(),
		CreatedAt: filter.CreatedAt,
		UpdatedAt: filter.UpdatedAt,
	}
}

func ModelListToFilterList(filters []models.Filter) []Filter {
	filterList := make([]Filter, 0, len(filters))
	for _, model := range filters {
		filterList = append(filterList, ModelToFilter(model))
	}
	return filterList
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
