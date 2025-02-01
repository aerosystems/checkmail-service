package usecases

import (
	"context"
	CustomErrors "github.com/aerosystems/checkmail-service/internal/common/custom_errors"
	"github.com/aerosystems/checkmail-service/internal/models"
	"time"
)

type AccessUsecase struct {
	apiAccessRepo AccessRepository
}

func NewAccessUsecase(apiAccessRepo AccessRepository) *AccessUsecase {
	return &AccessUsecase{apiAccessRepo: apiAccessRepo}
}

func (a AccessUsecase) GetAccess(ctx context.Context, token string) (*models.Access, error) {
	access, err := a.apiAccessRepo.Get(ctx, token)
	if err != nil {
		return nil, err
	}
	if access.AccessTime.Before(time.Now()) {
		return nil, CustomErrors.ErrSubscriptionIsNotActive
	}
	return access, nil
}

func (a AccessUsecase) CreateAccess(ctx context.Context, token, subscriptionType string, accessTime time.Time) error {
	return a.apiAccessRepo.CreateOrUpdate(ctx, &models.Access{
		Token:            token,
		SubscriptionType: models.SubscriptionTypeFromString(subscriptionType),
		AccessTime:       accessTime,
	})
}
