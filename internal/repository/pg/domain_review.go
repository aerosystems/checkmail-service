package pg

import (
	"github.com/aerosystems/checkmail-service/internal/models"
	"gorm.io/gorm"
)

type DomainReviewRepo struct {
	db *gorm.DB
}

func NewDomainReviewRepo(db *gorm.DB) *DomainReviewRepo {
	return &DomainReviewRepo{
		db: db,
	}
}

func (r *DomainReviewRepo) FindByName(name string) (*models.DomainReview, error) {
	var domainReview models.DomainReview
	result := r.db.First(&domainReview, "name = ?", name)
	if result.Error != nil {
		return nil, result.Error
	}
	return &domainReview, nil
}

func (r *DomainReviewRepo) Create(domainReview *models.DomainReview) error {
	result := r.db.Create(&domainReview)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func (r *DomainReviewRepo) Delete(domainReview *models.DomainReview) error {
	result := r.db.Delete(&domainReview)
	if result.Error != nil {
		return result.Error
	}
	return nil
}
