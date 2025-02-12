package usecases

import (
	"context"
	"errors"
	"github.com/aerosystems/checkmail-service/internal/entities"
)

type ManageUsecase struct {
	domainRepo DomainRepository
	filterRepo FilterRepository
}

func NewManageUsecase(domainRepo DomainRepository, filterRepo FilterRepository) *ManageUsecase {
	return &ManageUsecase{
		domainRepo: domainRepo,
		filterRepo: filterRepo,
	}
}

func (mu ManageUsecase) CreateDomain(ctx context.Context, domainName, domainType, domainCoverage string) (*entities.Domain, error) {
	domain := &entities.Domain{
		Name:  domainName,
		Type:  entities.DomainTypeFromString(domainType),
		Match: entities.DomainMatchFromString(domainCoverage),
	}
	if err := mu.domainRepo.Create(ctx, domain); err != nil {
		return nil, err // TODO: how to handle in handler http.StatusConflict or http.StatusInternalServerError?
	}
	return domain, nil
}

func (mu ManageUsecase) GetDomainByName(ctx context.Context, domainName string) (*entities.Domain, error) {
	return mu.domainRepo.FindByName(ctx, domainName)
}

func (mu ManageUsecase) UpdateDomain(ctx context.Context, domainName string, domainType, domainCoverage string) (*entities.Domain, error) {
	d, err := mu.domainRepo.FindByName(ctx, domainName)
	if err != nil {
		return nil, err
	}
	d.Type = entities.DomainTypeFromString(domainType)
	d.Match = entities.DomainMatchFromString(domainCoverage)
	if err := mu.domainRepo.Update(ctx, d); err != nil {
		return nil, err
	}
	return d, nil
}

func (mu ManageUsecase) DeleteDomain(ctx context.Context, domainName string) error {
	domain, err := mu.domainRepo.FindByName(ctx, domainName)
	if err != nil {
		return err
	}
	return mu.domainRepo.Delete(ctx, domain)
}

func (mu ManageUsecase) CountDomains(ctx context.Context) (map[entities.Type]int, error) {
	// TODO: add cache
	return mu.domainRepo.CountDomainTypes(ctx)
}

func (mu ManageUsecase) CreateFilter(ctx context.Context, domainName, domainType, domainCoverage, projectToken string) (entities.Filter, error) {
	return entities.Filter{}, errors.New("not implemented")
}
