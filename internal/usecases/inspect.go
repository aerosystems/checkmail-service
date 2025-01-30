package usecases

import (
	"context"
	CustomErrors "github.com/aerosystems/checkmail-service/internal/common/custom_errors"
	"github.com/aerosystems/checkmail-service/internal/common/validators"
	"github.com/aerosystems/checkmail-service/internal/models"
	"github.com/sirupsen/logrus"
	"golang.org/x/net/publicsuffix"
	"net/mail"
	"strings"
	"sync"
)

type InspectUsecase struct {
	log        *logrus.Logger
	domainRepo DomainRepository
	filterRepo FilterRepository
}

func NewInspectUsecase(
	log *logrus.Logger,
	domainRepo DomainRepository,
	filterRepo FilterRepository,
) *InspectUsecase {
	return &InspectUsecase{
		log:        log,
		domainRepo: domainRepo,
		filterRepo: filterRepo,
	}
}

func (i InspectUsecase) InspectData(ctx context.Context, data, clientIp, projectToken string) (models.Type, error) {
	domainName, err := getDomainName(data)
	if err != nil {
		return models.UndefinedType, CustomErrors.ErrEmailNotValid
	}

	if err = validators.ValidateDomainName(domainName); err != nil {
		return models.UndefinedType, err
	}

	if !isDomainExist(domainName) {
		return models.UndefinedType, CustomErrors.ErrDomainNotExist
	}

	return i.getDomainType(ctx, domainName)
}

func (i InspectUsecase) getDomainType(ctx context.Context, domainName string) (domainType models.Type, err error) {
	resChan := make(chan models.Type)
	errChan := make(chan error)

	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	wg := sync.WaitGroup{}
	wg.Add(2)
	go i.searchDomain(ctx, domainName, resChan, errChan, i.domainRepo.MatchEquals, i.domainRepo.MatchPrefix, i.domainRepo.MatchSuffix, i.domainRepo.MatchContains)

	domainType = models.UndefinedType
	go func() {
		defer wg.Done()
		for {
			select {
			case domainType = <-resChan:
				return
			case err = <-errChan:
				if err != nil {
					i.log.Errorf("error searching domain %s: %v", domainName, err)
				}
				return
			}
		}
	}()

	wg.Wait()
	return
}

func (i InspectUsecase) searchDomain(ctx context.Context, domainName string, resChan chan<- models.Type,
	errChan chan<- error, funcMatch ...func(ctx context.Context, name string) (*models.Domain, error)) {
	for _, f := range funcMatch {
		go func() {
			res, err := f(ctx, domainName)
			if err != nil {
				errChan <- err
				return
			}
			if res != nil {
				resChan <- res.Type
			}
		}()
	}
}

func getDomainName(data string) (string, error) {
	lowerData := strings.ToLower(data)
	if strings.Contains(lowerData, "@") {
		email, err := mail.ParseAddress(lowerData)
		if err != nil {
			return "", err
		}
		return strings.Split(email.Address, "@")[1], nil
	}
	return lowerData, nil
}

func isDomainExist(domainName string) bool {
	eTLD, icann := publicsuffix.PublicSuffix(domainName)
	if icann || strings.IndexByte(eTLD, '.') >= 0 {
		return true
	}
	return false
}
