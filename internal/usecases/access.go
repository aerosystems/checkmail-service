package usecases

import (
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
	access, err := a.apiAccessRepo.Get(apiKey)
	if err != nil {
		return nil, err
	}
	if access.AccessTime.Before(time.Now()) {
		return nil, CustomErrors.ErrSubscriptionIsNotActive
	}
	return access, nil
}
