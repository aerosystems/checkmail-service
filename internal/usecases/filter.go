package usecases

import (
	"errors"
	"github.com/aerosystems/checkmail-service/internal/models"
)

type FilterUsecase struct {
	rootDomainRepo RootDomainRepository
	filterRepo     FilterRepository
	projectRepo    ProjectRepository
}

func NewFilterUsecase(rootDomainRepo RootDomainRepository, filterRepo FilterRepository, projectRepo ProjectRepository) *FilterUsecase {
	return &FilterUsecase{
		rootDomainRepo: rootDomainRepo,
		filterRepo:     filterRepo,
		projectRepo:    projectRepo,
	}
}

func (fu *FilterUsecase) CreateFilter(domainName, domainType, domainCoverage, projectToken string) (models.Filter, error) {
	root, _ := GetRootDomain(domainName)
	rootDomain, _ := fu.rootDomainRepo.FindByName(root)
	if rootDomain != nil {
		err := errors.New("domain already exists")
		return models.Filter{}, err // http.StatusConflict
	}

	project, err := fu.projectRepo.GetProject(projectToken)
	if err != nil {
		return models.Filter{}, err
	}

	if project.Token != projectToken {
		err := errors.New("access denied")
		return models.Filter{}, err // http.StatusForbidden
	}

	filter := models.Filter{
		Name:         domainName,
		Type:         domainType,
		Coverage:     domainCoverage,
		ProjectToken: projectToken,
	}
	if err := fu.filterRepo.Create(&filter); err != nil {
		return models.Filter{}, err
	}
	return filter, nil
}
