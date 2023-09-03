package repository

import (
	"github.com/aerosystems/checkmail-service/internal/models"
	"gorm.io/gorm"
)

type FilterRepo struct {
	db *gorm.DB
}

func NewFilterRepo(db *gorm.DB) *FilterRepo {
	return &FilterRepo{
		db: db,
	}
}

func (r *FilterRepo) FindAll() (*[]models.Filter, error) {
	var filters []models.Filter
	result := r.db.Find(&filters)
	if result.Error != nil {
		return nil, result.Error
	}
	return &filters, nil
}

func (r *FilterRepo) FindByID(id int) (*models.Filter, error) {
	var filter models.Filter
	result := r.db.First(&filter, id)
	if result.Error != nil {
		return nil, result.Error
	}
	return &filter, nil
}

func (r *FilterRepo) FindByName(name string) (*models.Filter, error) {
	var filter models.Filter
	result := r.db.First(&filter, "name = ?", name)
	if result.Error != nil {
		return nil, result.Error
	}
	return &filter, nil
}

func (r *FilterRepo) Create(filter *models.Filter) error {
	result := r.db.Create(&filter)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func (r *FilterRepo) Update(filter *models.Filter) error {
	result := r.db.Save(&filter)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func (r *FilterRepo) Delete(filter *models.Filter) error {
	result := r.db.Delete(&filter)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func (r *FilterRepo) MatchEquals(name string) (*models.Filter, error) {
	var filter models.Filter
	result := r.db.First(&filter, "name = ? AND coverage = ?", name, "equals")
	if result.Error != nil {
		return nil, result.Error
	}
	return &filter, nil
}

func (r *FilterRepo) MatchContains(name string) (*models.Filter, error) {
	var filter models.Filter
	result := r.db.First(&filter, "name LIKE ? AND coverage = ?", "%"+name+"%", "contains")
	if result.Error != nil {
		return nil, result.Error
	}
	return &filter, nil
}

func (r *FilterRepo) MatchBegins(name string) (*models.Filter, error) {
	var filter models.Filter
	result := r.db.First(&filter, "name LIKE ? AND coverage = ?", name+"%", "begins")
	if result.Error != nil {
		return nil, result.Error
	}
	return &filter, nil
}

func (r *FilterRepo) MatchEnds(name string) (*models.Filter, error) {
	var filter models.Filter
	result := r.db.First(&filter, "name LIKE ? AND coverage = ?", "%"+name, "ends")
	if result.Error != nil {
		return nil, result.Error
	}
	return &filter, nil
}
