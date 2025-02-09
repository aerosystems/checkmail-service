package HTTPServer

import (
	"context"
	"github.com/aerosystems/checkmail-service/internal/models"
	"time"
)

type AccessUsecase interface {
	GetAccess(ctx context.Context, token string) (*models.Access, error)
	CreateAccess(ctx context.Context, token, subscriptionType string, accessCount int, accessTime time.Time) error
}

type InspectUsecase interface {
	InspectData(ctx context.Context, data, clientIp, projectToken string) (*models.Type, error)
	DeprecatedInspectData(ctx context.Context, data, clientIp, projectToken string) (*models.Type, error)
}

type ManageUsecase interface {
	CreateDomain(ctx context.Context, domainName, domainType, domainCoverage string) (*models.Domain, error)
	GetDomainByName(ctx context.Context, domainName string) (*models.Domain, error)
	UpdateDomain(ctx context.Context, domainName string, domainType, domainCoverage string) (*models.Domain, error)
	DeleteDomain(ctx context.Context, domainName string) error
	CountDomains(ctx context.Context) (map[models.Type]int, error)
	CreateFilter(ctx context.Context, domainName, domainType, domainCoverage, projectToken string) (models.Filter, error)
}

type ReviewUsecase interface {
	CreateReview(ctx context.Context, domainName, domainType string) (models.Review, error)
}
