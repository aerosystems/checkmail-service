package pg

import (
	"errors"
	"github.com/aerosystems/checkmail-service/internal/models"
	"gorm.io/gorm"
	"strings"
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

func (r *DomainRepo) FindById(id int) (*models.Domain, error) {
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
		if errors.Is(result.Error, gorm.ErrDuplicatedKey) || strings.Contains(result.Error.Error(), "duplicate key value violates unique constraint") {
			return errors.New("domain already exists")
		}
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

func (r *DomainRepo) MatchEquals(name string) (*models.Domain, error) {
	var domain models.Domain
	result := r.db.First(&domain, "name = ? AND coverage = ?", name, "equals")
	if result.Error != nil {
		return nil, result.Error
	}
	return &domain, nil
}

func (r *DomainRepo) MatchContains(name string) (*models.Domain, error) {
	var domain models.Domain
	result := r.db.First(&domain, "name LIKE ? AND coverage = ?", "%"+name+"%", "contains")
	if result.Error != nil {
		return nil, result.Error
	}
	return &domain, nil
}

func (r *DomainRepo) MatchBegins(name string) (*models.Domain, error) {
	var domain models.Domain
	result := r.db.First(&domain, "name LIKE ? AND coverage = ?", name+"%", "begins")
	if result.Error != nil {
		return nil, result.Error
	}
	return &domain, nil
}

func (r *DomainRepo) MatchEnds(name string) (*models.Domain, error) {
	var domain models.Domain
	result := r.db.First(&domain, "name LIKE ? AND coverage = ?", "%"+name, "ends")
	if result.Error != nil {
		return nil, result.Error
	}
	return &domain, nil
}

type TypeCount struct {
	Type  string
	Count int
}

func (r *DomainRepo) Count() (map[string]int, error) {
	var typeCounts []TypeCount
	if err := r.db.Model(&models.Domain{}).
		Select("type, COUNT(*) as count").
		Group("type").
		Scan(&typeCounts).Error; err != nil {
		return nil, err
	}

	var typeCountMap = make(map[string]int)
	for _, typeCount := range typeCounts {
		typeCountMap[typeCount.Type] = typeCount.Count
	}
	return typeCountMap, nil
}
