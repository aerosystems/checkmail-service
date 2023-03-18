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

func (r *DomainRepo) FindByID(ID int) (*models.Domain, error) {
	var project models.Domain
	result := r.db.First(&project, ID)
	if result.Error != nil {
		return nil, result.Error
	}
	return &project, nil
}

func (r *DomainRepo) FindByName(Token string) (*models.Domain, error) {
	var project models.Domain
	result := r.db.First(&project, "token = ?", Token)
	if result.Error != nil {
		return nil, result.Error
	}
	return &project, nil
}

func (r *DomainRepo) Create(project *models.Domain) error {
	result := r.db.Create(&project)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func (r *DomainRepo) Update(project *models.Domain) error {
	result := r.db.Save(&project)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func (r *DomainRepo) Delete(project *models.Domain) error {
	result := r.db.Delete(&project)
	if result.Error != nil {
		return result.Error
	}
	return nil
}
