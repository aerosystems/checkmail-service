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

func (a AccessUsecase) CreateAccess(token, subscriptionType string, accessTime time.Time) (*models.Access, error) {
	ctx := context.Background()
	access := models.Access{
		Token:            token,
		SubscriptionType: models.SubscriptionTypeFromString(subscriptionType),
		AccessTime:       accessTime,
	}
	if err := a.apiAccessRepo.Create(ctx, access); err != nil {
		return nil, err
	}
	return &access, nil
}
