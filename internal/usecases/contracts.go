package usecases

import (
	"context"
	"github.com/aerosystems/checkmail-service/internal/entities"
)

type DomainRepository interface {
	FindByName(ctx context.Context, name string) (*entities.Domain, error)
	Create(ctx context.Context, domain *entities.Domain) error
	Update(ctx context.Context, domain *entities.Domain) error
	Delete(ctx context.Context, domain *entities.Domain) error
	CountDomainTypes(ctx context.Context) (map[entities.Type]int, error)
	MatchEquals(ctx context.Context, name string) (*entities.Domain, error)
	MatchPrefix(ctx context.Context, name string) (*entities.Domain, error)
	MatchSuffix(ctx context.Context, name string) (*entities.Domain, error)
	MatchContains(ctx context.Context, name string) (*entities.Domain, error)
}

type FilterRepository interface {
	FindAll() ([]entities.Filter, error)
	FindByName(name string) (*entities.Filter, error)
	FindByProjectToken(projectToken string) ([]entities.Filter, error)
	Create(domain *entities.Filter) error
	CreateOrUpdate(domain *entities.Filter) error
	Delete(domain *entities.Filter) error
	MatchEquals(domainName, projectToken string) (*entities.Filter, error)
	MatchSuffix(domainName, projectToken string) (*entities.Filter, error)
}

type ReviewRepository interface {
	Create(domainReview *entities.Review) error
}

type AccessRepository interface {
	Get(ctx context.Context, token string) (*entities.Access, error)
	CreateOrUpdate(ctx context.Context, access *entities.Access) error
	Tx(ctx context.Context, token string, fn func(a *entities.Access) (any, error)) (any, error)
}

type LookupAdapter interface {
	Lookup(ctx context.Context, domain string) (entities.Type, error)
}
