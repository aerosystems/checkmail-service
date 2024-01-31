package rest

import (
	"github.com/aerosystems/checkmail-service/internal/models"
	CustomError "github.com/aerosystems/checkmail-service/pkg/custom_error"
)

type DomainRepository interface {
	FindByName(name string) (*models.Domain, error)
	Create(domain *models.Domain) error
	Update(domain *models.Domain) error
	Delete(domain *models.Domain) error
	Count() (map[string]int, error)
}

type RootDomainRepository interface {
	FindById(id int) (*models.RootDomain, error)
	FindByName(name string) (*models.RootDomain, error)
	Create(rootDomain *models.RootDomain) error
	Update(rootDomain *models.RootDomain) error
	Delete(rootDomain *models.RootDomain) error
}

type FilterRepository interface {
	FindAll() (*[]models.Filter, error)
	FindById(id int) (*models.Filter, error)
	FindByProjectToken(projectToken string) (*models.Filter, error)
	Create(domain *models.Filter) error
	Update(domain *models.Filter) error
	Delete(domain *models.Filter) error
}

type DomainReviewRepository interface {
	Create(domain *models.DomainReview) error
}

type InspectService interface {
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
}
