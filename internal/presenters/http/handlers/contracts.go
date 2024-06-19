package handlers

import (
	"github.com/aerosystems/checkmail-service/internal/models"
)

type InspectUsecase interface {
	InspectData(data, clientIp, projectToken string) (*string, *models.Error)
}

type DomainUsecase interface {
	CreateDomain(domainName, domainType, domainCoverage string) (*models.Domain, error)
	GetDomainByName(domainName string) (*models.Domain, error)
	UpdateDomain(domainName string, domainType, domainCoverage string) error
	DeleteDomain(domainName string) error
	CountDomains() (map[string]int, error)
}

type FilterUsecase interface {
	CreateFilter(domainName, domainType, domainCoverage, projectToken string) (models.Filter, error)
}

type ReviewUsecase interface {
	CreateReview(domainName, domainType string) (models.Review, error)
}
