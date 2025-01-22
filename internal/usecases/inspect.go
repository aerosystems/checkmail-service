package usecases

import (
	"context"
	"errors"
	CustomErrors "github.com/aerosystems/checkmail-service/internal/common/custom_errors"
	"github.com/aerosystems/checkmail-service/internal/common/validators"
	"github.com/aerosystems/checkmail-service/internal/models"
	"github.com/sirupsen/logrus"
	"net/mail"
	"net/rpc"
	"strings"
	"sync"
	"time"
)

type InspectUsecase struct {
	log            *logrus.Logger
	domainRepo     DomainRepository
	rootDomainRepo RootDomainRepository
	filterRepo     FilterRepository
}

type LookupRPCPayload struct {
	Domain   string
	ClientIp string
}

func NewInspectUsecase(
	log *logrus.Logger,
	domainRepo DomainRepository,
	rootDomainRepo RootDomainRepository,
	filterRepo FilterRepository,
) *InspectUsecase {
	return &InspectUsecase{
		log:            log,
		domainRepo:     domainRepo,
		rootDomainRepo: rootDomainRepo,
		filterRepo:     filterRepo,
	}
}

func (i *InspectUsecase) InspectData(data, clientIp, projectToken string) (models.Type, error) {
	start := time.Now()
	domainName, err := getDomainName(data)
	if err != nil {
		var publicApiError CustomErrors.PublicApiError
		if errors.As(err, &publicApiError) {
			i.setLogErrorRecord(data, publicApiError.Code, publicApiError.Message, projectToken, start)
		}
		return models.UndefinedType, CustomErrors.ErrEmailNotValid
	}

	if err = validators.ValidateDomainName(domainName); err != nil {
		var publicApiError CustomErrors.PublicApiError
		if errors.As(err, &publicApiError) {
			i.setLogErrorRecord(data, publicApiError.Code, publicApiError.Message, projectToken, start)
		}
		return models.UndefinedType, err
	}

	root := getRootDomainName(domainName)
	rootDomain, _ := i.rootDomainRepo.FindByName(root)
	if rootDomain == nil {
		var publicApiError CustomErrors.PublicApiError
		if errors.As(err, &publicApiError) {
			i.setLogErrorRecord(data, publicApiError.Code, publicApiError.Message, projectToken, start)
		}
		return models.UndefinedType, CustomErrors.ErrDomainNotExist
	}

	if projectToken != "" {
		if domainType := i.searchTypeDomainWithFilter(domainName, projectToken); domainType != models.UndefinedType {
			i.setLogSuccessRecord(data, domainName, domainType.String(), projectToken, "filter", start)
			return domainType, nil
		}
	}

	if domainType := i.searchTypeDomain(domainName); domainType != "undefined" {
		i.setLogSuccessRecord(data, domainName, domainType, projectToken, "database", start)
		return models.DomainTypeFromString(domainType), nil
	}

	var result string
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	errChan := make(chan error)

	go func(ctx context.Context) {
		lookupClientRPC, err := rpc.Dial("tcp", "lookup-service:5001")
		if err != nil {
			errChan <- err
		}

		errChan <- lookupClientRPC.Call("LookupServer.CheckDomain",
			LookupRPCPayload{Domain: domainName,
				ClientIp: clientIp},
			&result,
		)
	}(ctx)

	domainType := "undefined"
	select {
	case <-ctx.Done():
		i.log.Info("timeout check domain in lookup service via RPC")
		i.setLogSuccessRecord(data, domainName, domainType, projectToken, "", start)
		return models.DomainTypeFromString(domainType), nil
	case err := <-errChan:
		if err != nil {
			i.log.Warning("failed to check domain in lookup service via RPC: %v, result: %s", err, result)
			i.setLogSuccessRecord(data, domainName, domainType, projectToken, "", start)
			return models.DomainTypeFromString(domainType), nil
		}
		if err := validators.ValidateDomainTypes(result); err != nil {
			i.log.Warning("failed to check domain in lookup service via RPC: %v, result: %s", err, result)
			i.setLogSuccessRecord(data, domainName, domainType, projectToken, "", start)
			return models.DomainTypeFromString(domainType), nil
		}
		domain := &models.Domain{
			Name:  domainName,
			Type:  models.DomainTypeFromString(result),
			Match: models.EqualsMatch,
		}
		if err := i.domainRepo.Create(domain); err != nil {
			i.log.Error(err)
		}
		i.setLogSuccessRecord(data, domainName, domainType, projectToken, "lookup", start)
		return models.DomainTypeFromString(result), nil
	}

	return models.DomainTypeFromString(domainType), nil
}

func (i *InspectUsecase) setLogSuccessRecord(data, domainName, domainType, projectToken, source string, start time.Time) {
	duration := time.Since(start).Milliseconds()
	i.log.WithFields(logrus.Fields{"eventType": "inspect", "rawData": data, "domain": domainName, "domainType": domainType, "projectToken": projectToken, "duration": duration, "sourceInspect": source}).Info()
}

func (i *InspectUsecase) setLogErrorRecord(data string, errorCode int, errorMessage, projectToken string, start time.Time) {
	duration := time.Since(start).Milliseconds()
	i.log.WithFields(logrus.Fields{"eventType": "inspect", "rawData": data, "errorCode": errorCode, "errorMessage": errorMessage, "projectToken": projectToken, "duration": duration}).Info()
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

func getRootDomainName(domainName string) string {
	arrDomain := strings.Split(domainName, ".")
	return arrDomain[len(arrDomain)-1]
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

func (i *InspectUsecase) searchTypeDomain(domainName string) string {
	domainType := "undefined"

	chMatchEquals := make(chan string)
	chMatchEnds := make(chan string)
	chQuit := make(chan bool)

	go func() {
		res, _ := i.domainRepo.MatchEquals(domainName)
		if res != nil {
			chMatchEquals <- res.Type.String()
		}
		chQuit <- true
	}()

	go func() {
		res, _ := i.domainRepo.MatchSuffix(domainName)
		if res != nil {
			chMatchEnds <- res.Type.String()
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
