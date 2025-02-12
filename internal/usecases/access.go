package usecases

import (
	"context"
	"github.com/aerosystems/checkmail-service/internal/entities"
	"time"
)

type AccessUsecase struct {
	apiAccessRepo AccessRepository
}

func NewAccessUsecase(apiAccessRepo AccessRepository) *AccessUsecase {
	return &AccessUsecase{apiAccessRepo: apiAccessRepo}
}

func (a AccessUsecase) GetAccess(ctx context.Context, token string) (*entities.Access, error) {
	access, err := a.apiAccessRepo.Get(ctx, token)
	if err != nil {
		return nil, err
	}
	if access.AccessTime.Before(time.Now()) {
		return nil, entities.ErrSubscriptionIsNotActive
	}
	return access, nil
}

func (a AccessUsecase) CreateAccess(ctx context.Context, token, subscriptionType string, accessCount int, accessTime time.Time) error {
	return a.apiAccessRepo.CreateOrUpdate(ctx, &entities.Access{
		Token:            token,
		SubscriptionType: entities.SubscriptionTypeFromString(subscriptionType),
		AccessCount:      accessCount,
		AccessTime:       accessTime,
	})
}
