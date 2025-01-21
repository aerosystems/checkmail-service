package usecases

import (
	"errors"
	"github.com/aerosystems/checkmail-service/internal/models"
)

type FilterUsecase struct {
	rootDomainRepo RootDomainRepository
	filterRepo     FilterRepository
}

func NewFilterUsecase(rootDomainRepo RootDomainRepository, filterRepo FilterRepository) *FilterUsecase {
	return &FilterUsecase{
		rootDomainRepo: rootDomainRepo,
		filterRepo:     filterRepo,
	}
}

func (fu *FilterUsecase) CreateFilter(domainName, domainType, domainCoverage, projectToken string) (models.Filter, error) {
	return models.Filter{}, errors.New("not implemented")
}
