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

func (i *InspectUsecase) InspectData(ctx context.Context, data, clientIp, projectToken string) (models.Type, error) {
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

	return i.searchTypeDomainWithFilter(domainName, projectToken), nil
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

func (i *InspectUsecase) searchTypeDomainWithFilter(domainName, projectToken string) models.Type {
	domainType := models.UndefinedType

	chMatchEquals := make(chan models.Type)
	chMatchEnds := make(chan models.Type)
	chQuit := make(chan bool)

	go func() {
		res, _ := i.filterRepo.MatchEquals(domainName, projectToken)
		if res != nil {
			chMatchEquals <- res.Type
		}
		chQuit <- true
	}()

	go func() {
		res, _ := i.filterRepo.MatchSuffix(domainName, projectToken)
		if res != nil {
			chMatchEnds <- res.Type
		}
		chQuit <- true
	}()

	var wg sync.WaitGroup
	wg.Add(1)

	go func() {
		i := 0
		defer wg.Done()
		for {
			select {
			case domainType = <-chMatchEquals:
				return
			case domainType = <-chMatchEnds:
				return
			case <-chQuit:
				i++
				if i == 2 {
					return
				}
			}
		}
	}()

	wg.Wait()

	return domainType
}

func isDomainExist(domainName string) bool {
	eTLD, icann := publicsuffix.PublicSuffix(domainName)
	if icann || strings.IndexByte(eTLD, '.') >= 0 {
		return true
	}
	return false
}
