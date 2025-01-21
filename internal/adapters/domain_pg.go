package adapters

import (
	"errors"
	"fmt"
	CustomErrors "github.com/aerosystems/checkmail-service/internal/common/custom_errors"
	"github.com/aerosystems/checkmail-service/internal/models"
	"gorm.io/gorm"
	"strings"
	"time"
)

const (
	PrefixMatch   = "prefix"
	SuffixMatch   = "suffix"
	EqualsMatch   = "equals"
	ContainsMatch = "contains"
)

type Domain struct {
	Name      string    `gorm:"uniqueIndex:idx_name_type_match"`
	Type      string    `gorm:"uniqueIndex:idx_name_type_match;type:domain_type"`
	Match     string    `gorm:"uniqueIndex:idx_name_type_match;type:match_type"`
	CreatedAt time.Time `gorm:"autoCreateTime"`
	UpdatedAt time.Time `gorm:"autoUpdateTime"`
}

type DomainRepo struct {
	db *gorm.DB
}

func NewDomainRepo(db *gorm.DB) *DomainRepo {
	return &DomainRepo{
		db: db,
	}
}

func ModelToDomain(model *models.Domain) *Domain {
	return &Domain{
		Name:      model.Name,
		Type:      model.Type.String(),
		Match:     model.Match.String(),
		CreatedAt: model.CreatedAt,
		UpdatedAt: model.UpdatedAt,
	}
}

func DomainToModel(domain *Domain) *models.Domain {
	return &models.Domain{
		Name:      domain.Name,
		Type:      models.DomainTypeFromString(domain.Type),
		Match:     models.DomainCoverageFromString(domain.Match),
		CreatedAt: domain.CreatedAt,
		UpdatedAt: domain.UpdatedAt,
	}
}

func (r *DomainRepo) FindByName(name string) (*models.Domain, error) {
	var domain Domain
	result := r.db.First(&domain, "name = ?", name)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, CustomErrors.ErrDomainNotFound
		}
		return nil, fmt.Errorf("error finding domain by name: %w", result.Error)
	}
	return DomainToModel(&domain), nil
}

func (r *DomainRepo) Create(domain *models.Domain) error {
	result := r.db.Create(ModelToDomain(domain))
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrDuplicatedKey) || strings.Contains(result.Error.Error(), "duplicate key value violates unique constraint") {
			return CustomErrors.ErrDomainNotFound
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
	var domain Domain
	result := r.db.First(&domain, "name = ? AND coverage = ?", name, EqualsMatch)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, CustomErrors.ErrDomainNotFound
	}
	return DomainToModel(&domain), nil
}

func (r *DomainRepo) MatchContains(name string) (*models.Domain, error) {
	var domain Domain
	result := r.db.First(&domain, "name LIKE ? AND coverage = ?", "%"+name+"%", ContainsMatch)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, result.Error
	}
	return DomainToModel(&domain), nil
}

func (r *DomainRepo) MatchPrefix(name string) (*models.Domain, error) {
	var domain Domain
	result := r.db.First(&domain, "name LIKE ? AND coverage = ?", name+"%", PrefixMatch)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, result.Error
	}
	return DomainToModel(&domain), nil
}

func (r *DomainRepo) MatchSuffix(name string) (*models.Domain, error) {
	var domain Domain
	result := r.db.First(&domain, "name LIKE ? AND coverage = ?", "%"+name, SuffixMatch)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, result.Error
	}
	return DomainToModel(&domain), nil
}

func (r *DomainRepo) Count() (map[models.Type]int, error) {
	var typeCounts []struct {
		Type  string
		Count int
	}

	if err := r.db.Model(&Domain{}).
		Select("type, COUNT(*) as count").
		Group("type").
		Scan(&typeCounts).Error; err != nil {
		return nil, err
	}

	var typeCountMap = make(map[models.Type]int)
	for _, typeCount := range typeCounts {
		typeCountMap[models.DomainTypeFromString(typeCount.Type)] = typeCount.Count
	}
	return typeCountMap, nil
}
