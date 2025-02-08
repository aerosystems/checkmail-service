package usecases

import (
	"context"
	"errors"
	"fmt"
	"github.com/aerosystems/checkmail-service/internal/models"
	"github.com/sirupsen/logrus"
	"golang.org/x/net/publicsuffix"
	"net/mail"
	"strings"
	"sync"
	"time"
)

const domainInsertionTimeout = 5 * time.Second

type InspectUsecase struct {
	log           *logrus.Logger
	accessRepo    AccessRepository
	domainRepo    DomainRepository
	filterRepo    FilterRepository
	lookupService LookupAdapter
}

func NewInspectUsecase(log *logrus.Logger, accessRepo AccessRepository, domainRepo DomainRepository, filterRepo FilterRepository, lookupService LookupAdapter) *InspectUsecase {
	return &InspectUsecase{
		log:           log,
		accessRepo:    accessRepo,
		domainRepo:    domainRepo,
		filterRepo:    filterRepo,
		lookupService: lookupService,
	}
}

func (iu InspectUsecase) InspectData(ctx context.Context, data, _, projectToken string) (*models.Type, error) {
	res, err := iu.accessRepo.Tx(ctx, projectToken, func(a *models.Access) (any, error) {
		if err := validateAccess(a); err != nil {
			return nil, err
		}

		domainName, err := extractDomainName(data)
		if err != nil || !isValidDomain(domainName) {
			return models.UndefinedType, models.ErrDomainNotExist
		}

		domainType, err := iu.getDomainType(ctx, domainName)
		if err != nil {
			if errors.Is(err, models.ErrDomainNotExist) {
				return iu.lookupDomain(ctx, domainName)
			}
			return nil, err
		}

		if *domainType != models.UndefinedType {
			a.AccessCount--
		}
		return domainType, nil
	})
	if err != nil {
		return nil, err
	}

	domainType, ok := res.(*models.Type)
	if !ok {
		return nil, fmt.Errorf("unexpected result type %T", res)
	}

	return domainType, nil
}

func (iu InspectUsecase) lookupDomain(ctx context.Context, domainName string) (*models.Type, error) {
	domainType, err := iu.lookupService.Lookup(ctx, domainName)
	if err != nil {
		return nil, err
	}

	if domainType == models.UndefinedType {
		return &models.UndefinedType, nil
	}

	bgCtx := context.Background()
	go func() {
		newCtx, cancel := context.WithTimeout(bgCtx, domainInsertionTimeout)
		defer cancel()

		if err = iu.domainRepo.Create(newCtx, &models.Domain{
			Name:  domainName,
			Type:  domainType,
			Match: models.EqualsMatch,
		}); err != nil {
			iu.log.WithError(err).Errorf("failed to create domain %s", domainName)
		}
	}()
	return &domainType, nil
}

func validateAccess(a *models.Access) error {
	if a.AccessTime.Before(time.Now()) {
		return models.ErrAccessSubscriptionExpired
	}
	if a.SubscriptionType != models.BusinessSubscriptionType && a.AccessCount <= 0 {
		return models.ErrAccessLimitExceeded
	}
	return nil
}

func extractDomainName(data string) (string, error) {
	data = strings.ToLower(data)
	if strings.Contains(data, "@") {
		email, err := mail.ParseAddress(data)
		if err != nil {
			return "", err
		}
		return strings.Split(email.Address, "@")[1], nil
	}
	return data, nil
}

func isValidDomain(domainName string) bool {
	if err := models.ValidateDomainName(domainName); err != nil {
		return false
	}
	eTLD, icann := publicsuffix.PublicSuffix(domainName)
	return icann || strings.Contains(eTLD, ".")
}

func (iu InspectUsecase) getDomainType(ctx context.Context, domainName string) (*models.Type, error) {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	var wg sync.WaitGroup
	resChan := make(chan *models.Domain, 1)
	errChan := make(chan error, 1)
	matchFuncs := []func(ctx context.Context, name string) (*models.Domain, error){
		iu.domainRepo.MatchEquals,
		iu.domainRepo.MatchPrefix,
		iu.domainRepo.MatchSuffix,
		iu.domainRepo.MatchContains,
	}

	for _, f := range matchFuncs {
		wg.Add(1)
		go func(matchFunc func(ctx context.Context, name string) (*models.Domain, error)) {
			defer wg.Done()
			domain, err := matchFunc(ctx, domainName)
			resChan <- domain
			errChan <- err
		}(f)
	}

	go func() {
		wg.Wait()
		close(resChan)
		close(errChan)
	}()

	c := 1
	for {
		select {
		case domain := <-resChan:
			if domain != nil {
				return &domain.Type, nil
			}
		case err := <-errChan:
			if err != nil && !errors.Is(err, models.ErrDomainNotExist) {
				c++
				if c == len(matchFuncs) {
					return nil, models.ErrDomainNotExist
				}
			}
		case <-ctx.Done():
			return &models.UndefinedType, ctx.Err()
		}
	}
}
