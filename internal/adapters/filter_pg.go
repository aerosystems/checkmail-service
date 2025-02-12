package adapters

import (
	"errors"
	"fmt"
	"github.com/aerosystems/checkmail-service/internal/entities"
	"gorm.io/gorm"
	"strings"
)

type Filter struct {
	ProjectToken string
	Domain
}

type FilterRepo struct {
	db *gorm.DB
}

func NewFilterRepo(db *gorm.DB) *FilterRepo {
	if err := AutoMigrateGORM(db); err != nil {
		panic(fmt.Sprintf("failed to AutoMigrateGORM Filter model: %v", err))
	}
	return &FilterRepo{
		db: db,
	}
}

func ModelToFilter(model *entities.Filter) *Filter {
	return &Filter{
		ProjectToken: model.ProjectToken,
		Domain: Domain{
			Name:      model.Name,
			Type:      model.Type.String(),
			Match:     model.Match.String(),
			CreatedAt: model.CreatedAt,
			UpdatedAt: model.UpdatedAt,
		},
	}
}

func FilterToModel(filter *Filter) *entities.Filter {
	return &entities.Filter{
		ProjectToken: filter.ProjectToken,
		Domain: entities.Domain{
			Name:      filter.Name,
			Type:      entities.DomainTypeFromString(filter.Type),
			Match:     entities.DomainMatchFromString(filter.Match),
			CreatedAt: filter.CreatedAt,
			UpdatedAt: filter.UpdatedAt,
		},
	}
}

func (r *FilterRepo) FindAll() ([]entities.Filter, error) {
	var filters []Filter
	result := r.db.Find(&filters)
	if result.Error != nil {
		return nil, result.Error
	}
	var models []entities.Filter
	for _, filter := range filters {
		models = append(models, *FilterToModel(&filter))
	}
	return models, nil
}

func (r *FilterRepo) FindByName(name string) (*entities.Filter, error) {
	var filter Filter
	result := r.db.First(&filter, "name = ?", name)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, entities.ErrDomainNotFound
		}
		return nil, fmt.Errorf("error finding domain by name: %w", result.Error)
	}
	return FilterToModel(&filter), nil
}

func (r *FilterRepo) FindByProjectToken(projectToken string) ([]entities.Filter, error) {
	var filters []Filter
	result := r.db.Find(&filters, "project_token = ?", projectToken)
	if result.Error != nil {
		return nil, result.Error
	}
	var models []entities.Filter
	for _, filter := range filters {
		models = append(models, *FilterToModel(&filter))
	}
	return models, nil
}

func (r *FilterRepo) Create(filter *entities.Filter) error {
	result := r.db.Create(ModelToFilter(filter))
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrDuplicatedKey) || strings.Contains(result.Error.Error(), "duplicate key value violates unique constraint") {
			return entities.ErrDomainNotFound
		}
		return result.Error
	}
	return nil
}

func (r *FilterRepo) CreateOrUpdate(filter *entities.Filter) error {
	result := r.db.Save(ModelToFilter(filter))
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func (r *FilterRepo) Delete(filter *entities.Filter) error {
	result := r.db.Delete(&filter)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func (r *FilterRepo) MatchEquals(name, projectToken string) (*entities.Filter, error) {
	var filter Filter
	result := r.db.First(&filter, "project_token = ? AND name = ? AND match = ?", projectToken, name, EqualsMatch)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, entities.ErrDomainNotFound
	}
	return FilterToModel(&filter), nil
}

func (r *FilterRepo) MatchContains(name, projectToken string) (*entities.Filter, error) {
	var filter Filter
	result := r.db.First(&filter, "project_token = ? AND name LIKE ? AND match = ?", projectToken, "%"+name+"%", ContainsMatch)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, result.Error
	}
	return FilterToModel(&filter), nil
}

func (r *FilterRepo) MatchPrefix(name, projectToken string) (*entities.Filter, error) {
	var filter Filter
	result := r.db.First(&filter, "project_token = ? AND name LIKE ? AND match = ?", projectToken, name+"%", "begins")
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, result.Error
	}
	return FilterToModel(&filter), nil
}

func (r *FilterRepo) MatchSuffix(name, projectToken string) (*entities.Filter, error) {
	var filter Filter
	result := r.db.First(&filter, "project_token = ? AND name LIKE ? AND match = ?", projectToken, "%"+name, "ends")
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, result.Error
	}
	return FilterToModel(&filter), nil
}
