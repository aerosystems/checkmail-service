package usecases

import (
	"github.com/aerosystems/checkmail-service/internal/models"
)

type ReviewUsecase struct {
	reviewRepo ReviewRepository
}

func NewReviewUsecase(reviewRepo ReviewRepository) *ReviewUsecase {
	return &ReviewUsecase{
		reviewRepo: reviewRepo,
	}
}

func (ru ReviewUsecase) CreateReview(domainName, domainType string) (models.Review, error) {
	review := models.Review{
		Name: domainName,
		Type: domainType,
	}
	if err := ru.reviewRepo.Create(&review); err != nil {
		return models.Review{}, err // http.StatusInternalServerError
	}
	return review, nil
}
