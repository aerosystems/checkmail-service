package rest

import (
	"github.com/aerosystems/checkmail-service/internal/models"
	CustomError "github.com/aerosystems/checkmail-service/pkg/custom_error"
)

type InspectUsecase interface {
	InspectData(data, clientIp, projectToken string) (*string, *CustomError.Error)
}

type DomainUsecase interface {
	CreateDomain(domainName, domainType, domainCoverage string) (models.Domain, error)
	GetDomainByName(domainName string) (*models.Domain, error)
	UpdateDomain(domain *models.Domain, domainType, domainCoverage string) error
	DeleteDomain(domain *models.Domain) error
	CountDomains() (map[string]int, error)
}

type FilterUsecase interface {
	CreateFilter(domainName, domainType, domainCoverage, projectToken string) (models.Filter, error)
}

type ReviewUsecase interface {
	CreateReview(domainName, domainType string) (models.Review, error)
}
