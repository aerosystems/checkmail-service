package HTTPServer

import (
	"github.com/aerosystems/checkmail-service/internal/models"
	"time"
)

type FilterHandler struct {
	*BaseHandler
	manageUsecase ManageUsecase
}

func NewFilterHandler(
	baseHandler *BaseHandler,
	manageUsecase ManageUsecase,
) *FilterHandler {
	return &FilterHandler{
		BaseHandler:   baseHandler,
		manageUsecase: manageUsecase,
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
