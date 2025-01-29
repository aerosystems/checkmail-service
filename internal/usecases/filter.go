package usecases

import (
	"errors"
	"github.com/aerosystems/checkmail-service/internal/models"
)

type FilterUsecase struct {
	filterRepo FilterRepository
}

func NewFilterUsecase(filterRepo FilterRepository) *FilterUsecase {
	return &FilterUsecase{
		filterRepo: filterRepo,
	}
}

func (fu *FilterUsecase) CreateFilter(domainName, domainType, domainCoverage, projectToken string) (models.Filter, error) {
	return models.Filter{}, errors.New("not implemented")
}
