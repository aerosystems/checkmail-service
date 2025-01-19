package adapters

import (
	"github.com/aerosystems/checkmail-service/internal/models"
	"gorm.io/gorm"
)

type RootDomainRepo struct {
	db *gorm.DB
}

func NewRootDomainRepo(db *gorm.DB) *RootDomainRepo {
	return &RootDomainRepo{
		db: db,
	}
}

func (r *RootDomainRepo) FindById(id int) (*models.RootDomain, error) {
	var rootDomain models.RootDomain
	result := r.db.First(&rootDomain, id)
	if result.Error != nil {
		return nil, result.Error
	}
	return &rootDomain, nil
}

func (r *RootDomainRepo) FindByName(name string) (*models.RootDomain, error) {
	var rootDomain models.RootDomain
	result := r.db.First(&rootDomain, "name = ?", name)
	if result.Error != nil {
		return nil, result.Error
	}
	return &rootDomain, nil
}

func (r *RootDomainRepo) Create(rootDomain *models.RootDomain) error {
	result := r.db.Create(&rootDomain)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func (r *RootDomainRepo) Update(rootDomain *models.RootDomain) error {
	result := r.db.Save(&rootDomain)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func (r *RootDomainRepo) Delete(rootDomain *models.RootDomain) error {
	result := r.db.Delete(&rootDomain)
	if result.Error != nil {
		return result.Error
	}
	return nil
}
