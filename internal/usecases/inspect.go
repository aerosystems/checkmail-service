package usecases

import (
	"context"
	"errors"
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

func NewInspectUsecase(log *logrus.Logger, accessRepo AccessRepository, domainRepo DomainRepository, filterRepo FilterRepository) *InspectUsecase {
	return &InspectUsecase{
		log:        log,
		accessRepo: accessRepo,
		domainRepo: domainRepo,
		filterRepo: filterRepo,
	}
}

func (i *InspectUsecase) InspectData(ctx context.Context, data, _, projectToken string) (models.Type, error) {
	res, err := i.accessRepo.Tx(ctx, projectToken, func(a *models.Access) (any, error) {
		if err := validateAccess(a); err != nil {
			return models.UndefinedType, err
		}

		domainName, err := extractDomainName(data)
		if err != nil || !isValidDomain(domainName) {
			return models.UndefinedType, models.ErrDomainNotExist
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
	return res.(models.Type), nil
}

func validateAccess(a *models.Access) error {
	if a.AccessTime.Before(time.Now()) {
		return errors.New("access token expired")
	}
	if a.SubscriptionType != models.BusinessSubscriptionType && a.AccessCount <= 0 {
		return errors.New("access count exceeded")
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

func (i *InspectUsecase) getDomainType(ctx context.Context, domainName string) (models.Type, error) {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	var wg sync.WaitGroup
	resChan := make(chan models.Type, 1)
	errChan := make(chan error, 1)
	matchFuncs := []func(ctx context.Context, name string) (*models.Domain, error){
		i.domainRepo.MatchEquals,
		i.domainRepo.MatchPrefix,
		i.domainRepo.MatchSuffix,
		i.domainRepo.MatchContains,
	}

	for _, f := range matchFuncs {
		wg.Add(1)
		go func(matchFunc func(ctx context.Context, name string) (*models.Domain, error)) {
			defer wg.Done()
			if domain, err := matchFunc(ctx, domainName); err == nil && domain != nil {
				select {
				case resChan <- domain.Type:
					cancel()
				default:
				}
			}
		}(f)
	}

	go func() {
		wg.Wait()
		close(resChan)
		close(errChan)
	}()

	return <-resChan, nil
}
