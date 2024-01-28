package pg

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

func (r *FilterRepo) FindById(id int) (*models.Filter, error) {
	var filter models.Filter
	result := r.db.First(&filter, id)
	if result.Error != nil {
		return nil, result.Error
	}
	return &filter, nil
}

func (r *FilterRepo) FindByProjectToken(projectToken string) (*models.Filter, error) {
	var filter models.Filter
	result := r.db.First(&filter, "project_token = ?", projectToken)
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

func (r *FilterRepo) MatchEquals(domainName, projectToken string) (*models.Filter, error) {
	var filter models.Filter
	result := r.db.First(&filter, "name = ? AND project_token = ? AND coverage = ?", domainName, projectToken, "equals")
	if result.Error != nil {
		return nil, result.Error
	}
	return &filter, nil
}

func (r *FilterRepo) MatchContains(domainName, projectToken string) (*models.Filter, error) {
	var filter models.Filter
	result := r.db.First(&filter, "name LIKE ? AND project_token = ? AND coverage = ?", "%"+domainName+"%", projectToken, "contains")
	if result.Error != nil {
		return nil, result.Error
	}
	return &filter, nil
}

func (r *FilterRepo) MatchBegins(domainName, projectToken string) (*models.Filter, error) {
	var filter models.Filter
	result := r.db.First(&filter, "name LIKE ? AND project_token = ? AND coverage = ?", domainName+"%", projectToken, "begins")
	if result.Error != nil {
		return nil, result.Error
	}
	return &filter, nil
}

func (r *FilterRepo) MatchEnds(domainName, projectToken string) (*models.Filter, error) {
	var filter models.Filter
	result := r.db.First(&filter, "name LIKE ? AND project_token = ? AND coverage = ?", "%"+domainName, projectToken, "ends")
	if result.Error != nil {
		return nil, result.Error
	}
	return &filter, nil
}
