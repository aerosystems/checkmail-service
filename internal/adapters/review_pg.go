package adapters

import (
	"github.com/aerosystems/checkmail-service/internal/models"
	"gorm.io/gorm"
	"time"
)

type ReviewRepo struct {
	db *gorm.DB
}

func NewReviewRepo(db *gorm.DB) *ReviewRepo {
	return &ReviewRepo{
		db: db,
	}
}

type Review struct {
	Id        int       `gorm:"primaryKey;unique;autoIncrement"`
	Name      string    `gorm:"index:idx_name"`
	Type      string    `gorm:"<-"`
	CreatedAt time.Time `gorm:"autoCreateTime"`
	UpdatedAt time.Time `gorm:"autoUpdateTime"`
}

func (r *Review) DomainToModel(domain *Review) *models.Review {
	return &models.Review{
		Id:        domain.Id,
		Name:      domain.Name,
		Type:      domain.Type,
		CreatedAt: domain.CreatedAt,
		UpdatedAt: domain.UpdatedAt,
	}
}

func ReviewToDomain(model *models.Review) *Review {
	return &Review{
		Id:        model.Id,
		Name:      model.Name,
		Type:      model.Type,
		CreatedAt: model.CreatedAt,
		UpdatedAt: model.UpdatedAt,
	}
}

func (r *ReviewRepo) FindByName(name string) (*models.Review, error) {
	var domainReview Review
	result := r.db.First(&domainReview, "name = ?", name)
	if result.Error != nil {
		return nil, result.Error
	}
	return domainReview.DomainToModel(&domainReview), nil
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
