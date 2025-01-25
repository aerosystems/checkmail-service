package HTTPServer

import (
	"context"
	"github.com/aerosystems/checkmail-service/internal/models"
	"time"
)

type AccessUsecase interface {
	CreateAccess(token, subscriptionType string, accessTime time.Time) (*models.Access, error)
	GetAccess(ctx context.Context, token string) (*models.Access, error)
}

type InspectUsecase interface {
	InspectData(data, clientIp, projectToken string) (models.Type, error)
}

type DomainUsecase interface {
	CreateDomain(domainName, domainType, domainCoverage string) (*models.Domain, error)
	GetDomainByName(domainName string) (*models.Domain, error)
	UpdateDomain(domainName string, domainType, domainCoverage string) (*models.Domain, error)
	DeleteDomain(domainName string) error
	CountDomains() (map[models.Type]int, error)
}

type FilterUsecase interface {
	CreateFilter(domainName, domainType, domainCoverage, projectToken string) (models.Filter, error)
}

type ReviewUsecase interface {
	CreateReview(domainName, domainType string) (models.Review, error)
}
