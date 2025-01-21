package usecases

import (
	"context"
	"github.com/aerosystems/checkmail-service/internal/models"
)

type DomainRepository interface {
	FindByName(name string) (*models.Domain, error)
	Create(domain *models.Domain) error
	Update(domain *models.Domain) error
	Delete(domain *models.Domain) error
	Count() (map[models.Type]int, error)
	MatchEquals(name string) (*models.Domain, error)
	MatchSuffix(name string) (*models.Domain, error)
}

type RootDomainRepository interface {
	FindById(id int) (*models.RootDomain, error)
	FindByName(name string) (*models.RootDomain, error)
	Create(rootDomain *models.RootDomain) error
	Update(rootDomain *models.RootDomain) error
	Delete(rootDomain *models.RootDomain) error
}

type FilterRepository interface {
	FindAll() ([]models.Filter, error)
	FindByName(name string) (*models.Filter, error)
	FindByProjectToken(projectToken string) ([]models.Filter, error)
	Create(domain *models.Filter) error
	CreateOrUpdate(domain *models.Filter) error
	Delete(domain *models.Filter) error
	MatchEquals(domainName, projectToken string) (*models.Filter, error)
	MatchSuffix(domainName, projectToken string) (*models.Filter, error)
}

type ReviewRepository interface {
	Create(domainReview *models.Review) error
}

type ApiAccessRepository interface {
	Get(ctx context.Context, token string) (*models.Access, error)
	Create(ctx context.Context, access models.Access) error
}
