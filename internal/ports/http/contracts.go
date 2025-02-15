package HTTPServer

import (
	"context"
	"github.com/aerosystems/checkmail-service/internal/entities"
	"time"
)

type AccessUsecase interface {
	GetAccess(ctx context.Context, token string) (*entities.Access, error)
	CreateAccess(ctx context.Context, token, subscriptionType string, accessCount int, accessTime time.Time) error
}

type InspectUsecase interface {
	InspectData(ctx context.Context, data string) (*entities.Type, error)
	InspectDataWithAuth(ctx context.Context, data, clientIp, projectToken string) (*entities.Type, error)
	InspectDataDeprecated(ctx context.Context, data, clientIp, projectToken string) (*entities.Type, error)
}

type ManageUsecase interface {
	CreateDomain(ctx context.Context, domainName, domainType, domainCoverage string) (*entities.Domain, error)
	GetDomainByName(ctx context.Context, domainName string) (*entities.Domain, error)
	UpdateDomain(ctx context.Context, domainName string, domainType, domainCoverage string) (*entities.Domain, error)
	DeleteDomain(ctx context.Context, domainName string) error
	CountDomains(ctx context.Context) (map[entities.Type]int, error)
	CreateFilter(ctx context.Context, domainName, domainType, domainCoverage, projectToken string) (entities.Filter, error)
}

type ReviewUsecase interface {
	CreateReview(ctx context.Context, domainName, domainType string) (entities.Review, error)
}
