package adapters

import (
	"context"
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
		Match:     models.DomainMatchFromString(domain.Match),
		CreatedAt: domain.CreatedAt,
		UpdatedAt: domain.UpdatedAt,
	}
}

func (r *DomainRepo) FindByName(ctx context.Context, name string) (*models.Domain, error) {
	var domain Domain
	result := r.db.WithContext(ctx).First(&domain, "name = ?", name)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, CustomErrors.ErrDomainNotFound
		}
		return nil, fmt.Errorf("error finding domain by name: %w", result.Error)
	}
	return DomainToModel(&domain), nil
}

func (r *DomainRepo) Create(ctx context.Context, domain *models.Domain) error {
	result := r.db.WithContext(ctx).Create(ModelToDomain(domain))
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrDuplicatedKey) || strings.Contains(result.Error.Error(), "duplicate key value violates unique constraint") {
			return CustomErrors.ErrDomainAlreadyExists
		}
		return result.Error
	}
	return nil
}

func (r *DomainRepo) Update(ctx context.Context, domain *models.Domain) error {
	result := r.db.WithContext(ctx).Save(&domain)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func (r *DomainRepo) Delete(ctx context.Context, domain *models.Domain) error {
	result := r.db.WithContext(ctx).Delete(&domain)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func (r *DomainRepo) MatchEquals(ctx context.Context, name string) (*models.Domain, error) {
	var domain Domain
	result := r.db.WithContext(ctx).First(&domain, "name = ? AND match = ?", name, EqualsMatch)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, CustomErrors.ErrDomainNotFound
	}
	if result.Error != nil {
		return nil, result.Error
	}
	return DomainToModel(&domain), nil
}

func (r *DomainRepo) MatchContains(ctx context.Context, name string) (*models.Domain, error) {
	var domain Domain
	result := r.db.WithContext(ctx).First(&domain, "name LIKE ? AND match = ?", "%"+name+"%", ContainsMatch)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, CustomErrors.ErrDomainNotFound
	}
	if result.Error != nil {
		return nil, result.Error
	}
	return DomainToModel(&domain), nil
}

func (r *DomainRepo) MatchPrefix(ctx context.Context, name string) (*models.Domain, error) {
	var domain Domain
	result := r.db.WithContext(ctx).First(&domain, "name LIKE ? AND match = ?", name+"%", PrefixMatch)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, CustomErrors.ErrDomainNotFound
	}
	if result.Error != nil {
		return nil, result.Error
	}
	return DomainToModel(&domain), nil
}

func (r *DomainRepo) MatchSuffix(ctx context.Context, name string) (*models.Domain, error) {
	var domain Domain
	result := r.db.WithContext(ctx).First(&domain, "name LIKE ? AND match = ?", "%"+name, SuffixMatch)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, CustomErrors.ErrDomainNotFound
	}
	if result.Error != nil {
		return nil, result.Error
	}
	return DomainToModel(&domain), nil
}

func (r *DomainRepo) CountDomainTypes(ctx context.Context) (map[models.Type]int, error) {
	var typeCounts []struct {
		Type  string
		Count int
	}

	if err := r.db.WithContext(ctx).Model(&Domain{}).
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
