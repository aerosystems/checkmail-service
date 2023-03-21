package repository

import (
	"github.com/aerosystems/checkmail-service/internal/models"
	"gorm.io/gorm"
)

type DomainRepo struct {
	db *gorm.DB
}

func NewDomainRepo(db *gorm.DB) *DomainRepo {
	return &DomainRepo{
		db: db,
	}
}

func (r *DomainRepo) FindAll() (*[]models.Domain, error) {
	var domains []models.Domain
	result := r.db.Find(&domains)
	if result.Error != nil {
		return nil, result.Error
	}
	return &domains, nil
}

func (r *DomainRepo) FindByID(id int) (*models.Domain, error) {
	var domain models.Domain
	result := r.db.First(&domain, id)
	if result.Error != nil {
		return nil, result.Error
	}
	return &domain, nil
}

func (r *DomainRepo) FindByName(name string) (*models.Domain, error) {
	var domain models.Domain
	result := r.db.First(&domain, "name = ?", name)
	if result.Error != nil {
		return nil, result.Error
	}
	return &domain, nil
}

func (r *DomainRepo) Create(domain *models.Domain) error {
	result := r.db.Create(&domain)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func (r *DomainRepo) Update(domain *models.Domain) error {
	result := r.db.Save(&domain)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func (r *DomainRepo) Delete(domain *models.Domain) error {
	result := r.db.Delete(&domain)
	if result.Error != nil {
		return result.Error
	}
	return nil
}
