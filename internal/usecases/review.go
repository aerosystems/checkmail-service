package usecases

import (
	"context"
	"github.com/aerosystems/checkmail-service/internal/entities"
)

type ReviewUsecase struct {
	reviewRepo ReviewRepository
}

func NewReviewUsecase(reviewRepo ReviewRepository) *ReviewUsecase {
	return &ReviewUsecase{
		reviewRepo: reviewRepo,
	}
}

func (ru ReviewUsecase) CreateReview(ctx context.Context, domainName, domainType string) (entities.Review, error) {
	review := entities.Review{
		Name: domainName,
		Type: domainType,
	}
	if err := ru.reviewRepo.Create(&review); err != nil {
		return entities.Review{}, err // http.StatusInternalServerError
	}
	return review, nil
}
