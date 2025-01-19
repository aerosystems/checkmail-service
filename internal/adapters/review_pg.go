package adapters

import (
	"github.com/aerosystems/checkmail-service/internal/models"
	"gorm.io/gorm"
)

type ReviewRepo struct {
	db *gorm.DB
}

func NewReviewRepo(db *gorm.DB) *ReviewRepo {
	return &ReviewRepo{
		db: db,
	}
}

func (r *ReviewRepo) FindByName(name string) (*models.Review, error) {
	var domainReview models.Review
	result := r.db.First(&domainReview, "name = ?", name)
	if result.Error != nil {
		return nil, result.Error
	}
	return &domainReview, nil
}

func (r *ReviewRepo) Create(domainReview *models.Review) error {
	result := r.db.Create(&domainReview)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func (r *ReviewRepo) Delete(domainReview *models.Review) error {
	result := r.db.Delete(&domainReview)
	if result.Error != nil {
		return result.Error
	}
	return nil
}
