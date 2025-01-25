package adapters

import (
	"context"
	"github.com/aerosystems/checkmail-service/internal/models"
	"github.com/hashicorp/golang-lru/v2/expirable"
	"time"
)

const (
	apiKeyCacheTTL = 60 * time.Second
)

type CachedApiAccessRepo struct {
	lru           *expirable.LRU[string, *models.Access]
	apiAccessRepo *ApiAccessRepo
}

func NewCachedApiAccessRepo(apiAccessRepo *ApiAccessRepo) *CachedApiAccessRepo {
	return &CachedApiAccessRepo{
		lru:           expirable.NewLRU[string, *models.Access](0, nil, apiKeyCacheTTL),
		apiAccessRepo: apiAccessRepo,
	}
}

func (c CachedApiAccessRepo) Get(ctx context.Context, token string) (*models.Access, error) {
	if cachedAccess, ok := c.lru.Get(token); ok {
		return cachedAccess, nil
	}
	access, err := c.apiAccessRepo.Get(ctx, token)
	if err != nil {
		return nil, err
	}
	c.lru.Add(token, access)

	return access, nil
}

func (c CachedApiAccessRepo) Create(ctx context.Context, access models.Access) error {
	err := c.apiAccessRepo.Create(ctx, access)
	if err != nil {
		return err
	}
	_ = c.lru.Add(access.Token, &access)
	return nil
}
