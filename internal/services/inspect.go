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
}

type RPCLookupPayload struct {
	Domain   string
	ClientIp string
}

func NewInspectService(
	log *logrus.Logger,
	domainRepo models.DomainRepository,
	rootDomainRepo models.RootDomainRepository) *InspectService {
	return &InspectService{
		log:            log,
		domainRepo:     domainRepo,
		rootDomainRepo: rootDomainRepo,
	}
}

func (i *InspectService) InspectData(data, clientIp string) (*string, *CustomError.Error) {
	start := time.Now()
	lowerData := strings.ToLower(data)

	// Get Domain Name
	var domainName string
	if strings.Contains(lowerData, "@") {
		email, err := mail.ParseAddress(lowerData)
		if err != nil {
			return nil, CustomError.New(400001, "email address does not valid")
		}
		arr := strings.Split(email.Address, "@")
		domainName = arr[1]
	} else {
		domainName = lowerData
	}

	// Validate Domain Name
	if err := validators.ValidateDomain(domainName); err != nil {
		return nil, CustomError.New(400002, err.Error())
	}

	// Check Root Domain Name
	arrDomain := strings.Split(domainName, ".")
	root := arrDomain[len(arrDomain)-1]
	rootDomain, _ := i.rootDomainRepo.FindByName(root)
	if rootDomain == nil {
		return nil, CustomError.New(400003, "domain does not exist")
	}

	domainType := i.searchTypeDomain(domainName)

	if domainType == "undefined" {
		// check domain in lookup service via RPC
		var result string
		ctx, cancel := context.WithTimeout(context.Background(), time.Second)
		defer cancel()

		errChan := make(chan error, 1)

		go func(ctx context.Context) {
			lookupClientRPC, err := rpc.Dial("tcp", "lookup-service:5001")
			errChan <- err

			errChan <- lookupClientRPC.Call("LookupServer.CheckDomain",
				RPCLookupPayload{Domain: domainName,
					ClientIp: clientIp},
				&result,
			)
		}(ctx)

		select {
		case <-ctx.Done():
			// If 1 second timeout reached, send a partial response and continue waiting for result
			duration := time.Since(start)
			i.log.WithFields(logrus.Fields{
				"rawData":  data,
				"domain":   domainName,
				"type":     domainType,
				"duration": duration.Milliseconds(),
				"source":   "lookup",
			}).Info("successfully checked domain in lookup service via RPC")
			return &domainType, nil
		case err := <-errChan:
			if err == nil && validators.ValidateDomainTypes(result) == nil {
				domain := &models.Domain{
					Name:     domainName,
					Type:     result,
					Coverage: "equals",
				}
				err := i.domainRepo.Create(domain)
				if err != nil {
					i.log.Error(err)
				}
			} else {
				i.log.Error(err)
			}
		}
	}

	duration := time.Since(start)
	i.log.WithFields(logrus.Fields{
		"rawData":  data,
		"domain":   domainName,
		"type":     domainType,
		"duration": duration.Milliseconds(),
		"source":   "database",
	}).Info("successfully checked domain in local database")
	return &domainType, nil
}

func (i *InspectService) searchTypeDomain(domainName string) string {

	domainType := "undefined"

	chMatchEquals := make(chan string)
	chMatchContains := make(chan string)
	chMatchBegins := make(chan string)
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
		res, _ := i.domainRepo.MatchContains(domainName)
		if res != nil {
			chMatchContains <- res.Type
		}
		chQuit <- true
	}()

	go func() {
		res, _ := i.domainRepo.MatchBegins(domainName)
		if res != nil {
			chMatchBegins <- res.Type
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
			case domainType = <-chMatchContains:
				return
			case domainType = <-chMatchBegins:
				return
			case domainType = <-chMatchEnds:
				return
			case <-chQuit:
				i++
				if i == 4 {
					return
				}
			}
		}
	}()

	wg.Wait()

	return domainType
}
