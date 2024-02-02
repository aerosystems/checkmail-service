package usecases

import (
	"errors"
	"github.com/aerosystems/checkmail-service/internal/models"
)

type ReviewUsecase struct {
	reviewRepo     ReviewRepository
	rootDomainRepo RootDomainRepository
}

func NewReviewUsecase(reviewRepo ReviewRepository, rootDomainRepo RootDomainRepository) *ReviewUsecase {
	return &ReviewUsecase{
		reviewRepo:     reviewRepo,
		rootDomainRepo: rootDomainRepo,
	}
}

func (ru ReviewUsecase) CreateReview(domainName, domainType string) (models.Review, error) {
	root, _ := GetRootDomain(domainName)
	rootDomain, _ := ru.rootDomainRepo.FindByName(root)
	if rootDomain == nil {
		err := errors.New("domain does not exist")
		return models.Review{}, err // http.StatusNotFound
	}
	review := models.Review{
		Name: domainName,
		Type: domainType,
	}
	if err := ru.reviewRepo.Create(&review); err != nil {
		return models.Review{}, err // http.StatusInternalServerError
	}
	return review, nil
}
