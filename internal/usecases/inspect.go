package usecases

import (
	"context"
	"errors"
	"fmt"
	"github.com/aerosystems/checkmail-service/internal/entities"
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

func (iu InspectUsecase) InspectData(ctx context.Context, data, _, projectToken string) (*entities.Type, error) {
	res, err := iu.accessRepo.Tx(ctx, projectToken, func(a *entities.Access) (any, error) {
		if err := validateAccess(a); err != nil {
			return nil, err
		}

		domainName, err := extractDomainName(data)
		if err != nil || !isValidDomain(domainName) {
			return entities.UndefinedType, entities.ErrDomainNotExist
		}

		domainType, err := iu.getDomainType(ctx, domainName)
		if err != nil {
			if errors.Is(err, entities.ErrDomainNotExist) {
				return iu.lookupDomain(ctx, domainName)
			}
			return nil, err
		}

		if *domainType != entities.UndefinedType {
			a.AccessCount--
		}
		return domainType, nil
	})
	if err != nil {
		return nil, err
	}

	domainType, ok := res.(*entities.Type)
	if !ok {
		return nil, fmt.Errorf("unexpected result type %T", res)
	}

	return domainType, nil
}

func (iu InspectUsecase) lookupDomain(ctx context.Context, domainName string) (*entities.Type, error) {
	domainType, err := iu.lookupService.Lookup(ctx, domainName)
	if err != nil {
		return nil, err
	}

	if domainType == entities.UndefinedType {
		return &entities.UndefinedType, nil
	}

	bgCtx := context.Background()
	go func() {
		newCtx, cancel := context.WithTimeout(bgCtx, domainInsertionTimeout)
		defer cancel()

		if err = iu.domainRepo.Create(newCtx, &entities.Domain{
			Name:  domainName,
			Type:  domainType,
			Match: entities.EqualsMatch,
		}); err != nil {
			iu.log.WithError(err).Errorf("failed to create domain %s", domainName)
		}
	}()
	return &domainType, nil
}

func validateAccess(a *entities.Access) error {
	if a.AccessTime.Before(time.Now()) {
		return entities.ErrAccessSubscriptionExpired
	}
	if a.SubscriptionType != entities.BusinessSubscriptionType && a.AccessCount <= 0 {
		return entities.ErrAccessLimitExceeded
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
	if err := entities.ValidateDomainName(domainName); err != nil {
		return false
	}
	eTLD, icann := publicsuffix.PublicSuffix(domainName)
	return icann || strings.Contains(eTLD, ".")
}

func (iu InspectUsecase) getDomainType(ctx context.Context, domainName string) (*entities.Type, error) {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	var wg sync.WaitGroup
	resChan := make(chan *entities.Domain, 1)
	errChan := make(chan error, 1)
	matchFuncs := []func(ctx context.Context, name string) (*entities.Domain, error){
		iu.domainRepo.MatchEquals,
		iu.domainRepo.MatchPrefix,
		iu.domainRepo.MatchSuffix,
		iu.domainRepo.MatchContains,
	}

	for _, f := range matchFuncs {
		wg.Add(1)
		go func(matchFunc func(ctx context.Context, name string) (*entities.Domain, error)) {
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
			if err != nil && !errors.Is(err, entities.ErrDomainNotExist) {
				c++
				if c == len(matchFuncs) {
					return nil, entities.ErrDomainNotExist
				}
			}
		case <-ctx.Done():
			return &entities.UndefinedType, ctx.Err()
		}
	}
}

func (iu InspectUsecase) DeprecatedInspectData(ctx context.Context, data, _, projectToken string) (*entities.Type, error) {
	return nil, nil
}
