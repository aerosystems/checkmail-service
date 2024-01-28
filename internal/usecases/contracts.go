package usecases

import "github.com/aerosystems/checkmail-service/internal/models"

type DomainRepository interface {
	Create(domain *models.Domain) error
	MatchEquals(name string) (*models.Domain, error)
	MatchEnds(name string) (*models.Domain, error)
}

type RootDomainRepository interface {
	FindByName(name string) (*models.RootDomain, error)
}

type FilterRepository interface {
	MatchEquals(domainName, projectToken string) (*models.Filter, error)
	MatchEnds(domainName, projectToken string) (*models.Filter, error)
}
