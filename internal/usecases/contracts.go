package usecases

import (
	"context"
	"github.com/aerosystems/checkmail-service/internal/models"
)

type DomainRepository interface {
	FindByName(ctx context.Context, name string) (*models.Domain, error)
	Create(ctx context.Context, domain *models.Domain) error
	Update(ctx context.Context, domain *models.Domain) error
	Delete(ctx context.Context, domain *models.Domain) error
	CountDomainTypes(ctx context.Context) (map[models.Type]int, error)
	MatchEquals(ctx context.Context, name string) (*models.Domain, error)
	MatchPrefix(ctx context.Context, name string) (*models.Domain, error)
	MatchSuffix(ctx context.Context, name string) (*models.Domain, error)
	MatchContains(ctx context.Context, name string) (*models.Domain, error)
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
