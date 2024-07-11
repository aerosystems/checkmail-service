package usecases

import (
	"context"
	CustomErrors "github.com/aerosystems/checkmail-service/internal/common/custom_errors"
	"github.com/aerosystems/checkmail-service/internal/models"
	"time"
)

type AccessUsecase struct {
	apiAccessRepo ApiAccessRepository
}

func NewAccessUsecase(apiAccessRepo ApiAccessRepository) *AccessUsecase {
	return &AccessUsecase{apiAccessRepo: apiAccessRepo}
}

func (a AccessUsecase) GetAccess(apiKey string) (*models.Access, error) {
	ctx := context.Background()
	access, err := a.apiAccessRepo.Get(ctx, apiKey)
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
