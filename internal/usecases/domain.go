package usecases

import (
	"errors"
	"github.com/aerosystems/checkmail-service/internal/models"
	"strings"
)

type DomainUsecase struct {
	domainRepo     DomainRepository
	rootDomainRepo RootDomainRepository
}

func NewDomainUsecase(domainRepo DomainRepository, rootDomainRepo RootDomainRepository) *DomainUsecase {
	return &DomainUsecase{
		domainRepo:     domainRepo,
		rootDomainRepo: rootDomainRepo,
	}
}

func (du *DomainUsecase) CreateDomain(domainName, domainType, domainCoverage string) (*models.Domain, error) {
	root, _ := GetRootDomain(domainName)
	rootDomain, _ := du.rootDomainRepo.FindByName(root)
	if rootDomain == nil {
		return nil, errors.New("domain does not exist") // http.StatusNotFound
	}
	domain := &models.Domain{
		Name:     domainName,
		Type:     models.DomainTypeFromString(domainType),
		Coverage: models.DomainCoverageFromString(domainCoverage),
	}
	if err := du.domainRepo.Create(domain); err != nil {
		return nil, err // TODO: how to handle in handler http.StatusConflict or http.StatusInternalServerError?
	}
	return domain, nil
}

func (du *DomainUsecase) GetDomainByName(domainName string) (*models.Domain, error) {
	return du.domainRepo.FindByName(domainName)
}

func (du *DomainUsecase) UpdateDomain(domainName string, domainType, domainCoverage string) error {
	domain, err := du.domainRepo.FindByName(domainName)
	if err != nil {
		return err
	}
	domain.Type = models.DomainTypeFromString(domainType)
	domain.Coverage = models.DomainCoverageFromString(domainCoverage)
	return du.domainRepo.Update(domain)
}

func (du *DomainUsecase) DeleteDomain(domainName string) error {
	domain, err := du.domainRepo.FindByName(domainName)
	if err != nil {
		return err
	}
	return du.domainRepo.Delete(domain)
}

func (du *DomainUsecase) CountDomains() (map[string]int, error) {
	// TODO: add cache
	return du.domainRepo.Count()
}

func GetRootDomain(domain string) (string, error) {
	arrDomain := strings.Split(domain, ".")
	if len(arrDomain) < 2 {
		return "", errors.New("domain is not valid")
	}
	return arrDomain[len(arrDomain)-1], nil
}
