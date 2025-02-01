package usecases

import (
	"context"
	"errors"
	CustomErrors "github.com/aerosystems/checkmail-service/internal/common/custom_errors"
	"github.com/aerosystems/checkmail-service/internal/common/validators"
	"github.com/aerosystems/checkmail-service/internal/models"
	"github.com/sirupsen/logrus"
	"golang.org/x/net/publicsuffix"
	"net/mail"
	"strings"
	"sync"
	"time"
)

type InspectUsecase struct {
	log        *logrus.Logger
	accessRepo AccessRepository
	domainRepo DomainRepository
	filterRepo FilterRepository
}

func NewInspectUsecase(
	log *logrus.Logger,
	accessRepo AccessRepository,
	domainRepo DomainRepository,
	filterRepo FilterRepository,
) *InspectUsecase {
	return &InspectUsecase{
		log:        log,
		accessRepo: accessRepo,
		domainRepo: domainRepo,
		filterRepo: filterRepo,
	}
}

func (i InspectUsecase) InspectData(ctx context.Context, data, _, projectToken string) (models.Type, error) {
	res, err := i.accessRepo.Tx(ctx, projectToken, func(a *models.Access) (any, error) {
		if a.AccessTime.Before(time.Now().Add(-time.Hour)) {
			return models.UndefinedType, errors.New("access token expired")
		}

		if a.AccessCount <= 0 {
			return models.UndefinedType, errors.New("access count exceeded")
		}

		domainName, err := getDomainName(data)
		if err != nil {
			return models.UndefinedType, err
		}

		if err = validators.ValidateDomainName(domainName); err != nil {
			return models.UndefinedType, err
		}

		if !isDomainExist(domainName) {
			return models.UndefinedType, CustomErrors.ErrDomainNotExist
		}

		domainType, err := i.getDomainType(ctx, domainName)
		if err != nil {
			return models.UndefinedType, err
		}
		if domainType != models.UndefinedType {
			a.AccessCount--
		}

		return domainType, nil
	})
	if err != nil {
		return models.UndefinedType, err
	}
	t, ok := res.(models.Type)
	if !ok {
		return models.UndefinedType, errors.New("unexpected result type")
	}
	return t, nil
}

func (i InspectUsecase) getDomainType(ctx context.Context, domainName string) (domainType models.Type, err error) {
	resChan := make(chan models.Type)
	errChan := make(chan error)

	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	matchFuncs := []func(ctx context.Context, name string) (*models.Domain, error){
		i.domainRepo.MatchEquals,
		i.domainRepo.MatchPrefix,
		i.domainRepo.MatchSuffix,
		i.domainRepo.MatchContains,
	}

	wg := sync.WaitGroup{}
	wg.Add(2)
	go func() {
		defer wg.Done()
		i.searchDomain(ctx, domainName, resChan, errChan, matchFuncs...)
	}()

	go func() {
		defer wg.Done()
		count := 0
		for {
			select {
			case domainType = <-resChan:
				if domainType != models.UndefinedType {
					cancel()
					return
				}
				count++
				if count >= len(matchFuncs) {
					cancel()
					return
				}
			case <-errChan:
				return
			}
		}
	}()

	wg.Wait()
	return domainType, err
}

func (i InspectUsecase) searchDomain(ctx context.Context, domainName string, resChan chan<- models.Type,
	errChan chan<- error, funcMatch ...func(ctx context.Context, name string) (*models.Domain, error)) {
	for _, f := range funcMatch {
		go func() {
			res, err := f(ctx, domainName)
			if err != nil && !errors.Is(err, CustomErrors.ErrDomainNotFound) {
				errChan <- err
				i.log.WithError(err).Error("search domain error")
			}
			if res != nil {
				resChan <- res.Type
			}
			resChan <- models.UndefinedType
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
