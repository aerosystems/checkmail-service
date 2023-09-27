package services

import (
	"context"
	"github.com/aerosystems/checkmail-service/internal/models"
	CustomError "github.com/aerosystems/checkmail-service/pkg/custom_error"
	"github.com/aerosystems/checkmail-service/pkg/validators"
	"github.com/sirupsen/logrus"
	"net/mail"
	"net/rpc"
	"strings"
	"sync"
	"time"
)

type InspectService struct {
	log            *logrus.Logger
	domainRepo     models.DomainRepository
	rootDomainRepo models.RootDomainRepository
	filterRepo     models.FilterRepository
}

type LookupRPCPayload struct {
	Domain   string
	ClientIp string
}

func NewInspectService(
	log *logrus.Logger,
	domainRepo models.DomainRepository,
	rootDomainRepo models.RootDomainRepository,
	filterRepo models.FilterRepository,
) *InspectService {
	return &InspectService{
		log:            log,
		domainRepo:     domainRepo,
		rootDomainRepo: rootDomainRepo,
		filterRepo:     filterRepo,
	}
}

func (i *InspectService) InspectData(data, clientIp, projectToken string) (*string, *CustomError.Error) {
	start := time.Now()
	domainName, err := getDomainName(data)
	if err != nil {
		return nil, CustomError.New(400001, "email address does not valid")
	}

	if err := validators.ValidateDomainName(domainName); err != nil {
		return nil, err
	}

	root := getRootDomainName(domainName)
	rootDomain, _ := i.rootDomainRepo.FindByName(root)
	if rootDomain == nil {
		return nil, CustomError.New(400003, "domain does not exist")
	}

	if projectToken != "" {
		if domainType := i.searchTypeDomainWithFilter(domainName, projectToken); domainType != "undefined" {
			i.setLogRecord(data, domainName, domainType, projectToken, "filter", start)
			return &domainType, nil
		}
	}

	if domainType := i.searchTypeDomain(domainName); domainType != "undefined" {
		i.setLogRecord(data, domainName, domainType, projectToken, "database", start)
		return &domainType, nil
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
		i.setLogRecord(data, domainName, domainType, projectToken, "", start)
		return &domainType, nil
	case err := <-errChan:
		if err != nil {
			i.log.Warning("failed to check domain in lookup service via RPC: %v, result: %s", err, result)
			i.setLogRecord(data, domainName, domainType, projectToken, "", start)
			return &domainType, nil
		}
		if err := validators.ValidateDomainTypes(result); err != nil {
			i.log.Warning("failed to check domain in lookup service via RPC: %v, result: %s", err, result)
			i.setLogRecord(data, domainName, domainType, projectToken, "", start)
			return &domainType, nil
		}
		domain := &models.Domain{
			Name:     domainName,
			Type:     result,
			Coverage: "equals",
		}
		if err := i.domainRepo.Create(domain); err != nil {
			i.log.Error(err)
		}
		i.setLogRecord(data, domainName, domainType, projectToken, "lookup", start)
		return &result, nil
	}

	return &domainType, nil
}

func (i *InspectService) setLogRecord(data, domainName, domainType, projectToken, source string, start time.Time) {
	duration := time.Since(start).Milliseconds()
	i.log.WithFields(logrus.Fields{"event": "inspect", "rawData": data, "domain": domainName, "type": domainType, "projectToken": projectToken, "duration": duration, "source": source}).Info()
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

func (i *InspectService) searchTypeDomainWithFilter(domainName, projectToken string) string {
	domainType := "undefined"

	chMatchEquals := make(chan string)
	chMatchEnds := make(chan string)
	chQuit := make(chan bool)

	go func() {
		res, _ := i.filterRepo.MatchEquals(domainName, projectToken)
		if res != nil {
			chMatchEquals <- res.Type
		}
		chQuit <- true
	}()

	go func() {
		res, _ := i.filterRepo.MatchEnds(domainName, projectToken)
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

func (i *InspectService) searchTypeDomain(domainName string) string {
	domainType := "undefined"

	chMatchEquals := make(chan string)
	chMatchEnds := make(chan string)
	chQuit := make(chan bool)

	go func() {
		res, _ := i.domainRepo.MatchEquals(domainName)
		if res != nil {
			chMatchEquals <- res.Type
		}
		chQuit <- true
	}()

	go func() {
		res, _ := i.domainRepo.MatchEnds(domainName)
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
