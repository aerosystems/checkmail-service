package adapters

import (
	"context"
	"errors"
	"github.com/aerosystems/checkmail-service/internal/entities"
	"gorm.io/gorm"
	"time"
)

type AccessRepo struct {
	db *gorm.DB
}

func NewAccessRepo(db *gorm.DB) *AccessRepo {
	return &AccessRepo{db: db}
}

type Access struct {
	Token            string    `gorm:"uniqueIndex:idx_token_access"`
	SubscriptionType string    `gorm:"uniqueIndex:idx_token_access;type:subscription_type"`
	AccessCount      int       `gorm:"uniqueIndex:idx_token_access"`
	AccessTime       time.Time `gorm:"<-"`
}

func (a *Access) ToModel() *entities.Access {
	return &entities.Access{
		Token:            a.Token,
		SubscriptionType: entities.SubscriptionTypeFromString(a.SubscriptionType),
		AccessCount:      a.AccessCount,
		AccessTime:       a.AccessTime,
	}
}

func ModelToAccess(access *entities.Access) *Access {
	return &Access{
		Token:            access.Token,
		SubscriptionType: access.SubscriptionType.String(),
		AccessCount:      access.AccessCount,
		AccessTime:       access.AccessTime,
	}
}

func (ar *AccessRepo) Get(ctx context.Context, token string) (*entities.Access, error) {
	var access Access
	result := ar.db.WithContext(ctx).Where("token = ?", token).First(&access)
	if result.Error != nil {
		return nil, result.Error
	}
	return access.ToModel(), nil
}

func (ar *AccessRepo) CreateOrUpdate(ctx context.Context, access *entities.Access) error {
	accessModel := ModelToAccess(access)
	result := ar.db.WithContext(ctx).Where("token = ?", accessModel.Token).First(&Access{})
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			result = ar.db.WithContext(ctx).Create(accessModel)
			if result.Error != nil {
				return result.Error
			}
			return nil
		}
		return result.Error
	}
	result = ar.db.WithContext(ctx).Model(&Access{}).Where("token = ?", accessModel.Token).Updates(accessModel)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func (ar *AccessRepo) Tx(ctx context.Context, token string, fn func(a *entities.Access) (any, error)) (any, error) {
	tx := ar.db.WithContext(ctx).Begin()
	if tx.Error != nil {
		return nil, tx.Error
	}
	var access Access
	err := tx.Model(&Access{}).Where("token = ?", token).First(&access).Error
	if err != nil {
		tx.Rollback()
		return nil, err
	}
	accessModel := access.ToModel()
	result, err := fn(accessModel)
	if err != nil {
		tx.Rollback()
		return nil, err
	}
	tx = tx.Model(&Access{}).Where("token = ?", token).Select("*").Updates(ModelToAccess(accessModel))
	if tx.Error != nil {
		tx.Rollback()
		return nil, tx.Error
	}
	return result, tx.Commit().Error
}
