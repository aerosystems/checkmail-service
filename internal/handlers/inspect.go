package handlers

import (
	"context"
	"errors"
	"fmt"
	"github.com/aerosystems/checkmail-service/internal/models"
	"github.com/aerosystems/checkmail-service/pkg/validators"
	"github.com/sirupsen/logrus"
	"net/http"
	"net/mail"
	"net/rpc"
	"strings"
	"sync"
	"time"
)

type RPCLookupPayload struct {
	Domain   string
	ClientIp string
}

type InspectRequestPayload struct {
	Data     string `json:"data"`
	ClientIp string `json:"clientIp,omitempty"`
}

// Inspect godoc
// @Summary get information about domain name or email address
// @Tags inspect
// @Accept  json
// @Produce application/json
// @Param data body InspectRequestPayload true "raw request body"
// @Security X-API-KEY
// @Success 200 {object} Response{data=string}
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /v1/inspect [post]
func (h *BaseHandler) Inspect(w http.ResponseWriter, r *http.Request) {
	start := time.Now()
	var requestPayload InspectRequestPayload
	if err := ReadRequest(w, r, &requestPayload); err != nil {
		_ = WriteResponse(w, http.StatusUnprocessableEntity, NewErrorPayload(422001, "could not read request body", err))
		return
	}

	data := strings.ToLower(requestPayload.Data)

	// Get Domain Name
	var domainName string
	if strings.Contains(data, "@") {
		email, err := mail.ParseAddress(data)
		if err != nil {
			err := errors.New("email address does not valid")
			_ = WriteResponse(w, http.StatusBadRequest, NewErrorPayload(400001, err.Error(), err))
			return
		}
		arr := strings.Split(email.Address, "@")
		domainName = arr[1]
	} else {
		domainName = data
	}

	// Validate Domain Name
	isValid := validators.ValidateDomain(domainName)
	if !isValid {
		err := errors.New("domain does not valid")
		_ = WriteResponse(w, http.StatusBadRequest, NewErrorPayload(400002, err.Error(), err))
		return
	}

	// Check Root Domain Name
	arrDomain := strings.Split(domainName, ".")
	root := arrDomain[len(arrDomain)-1]
	rootDomain, _ := h.rootDomainRepo.FindByName(root)
	if rootDomain == nil {
		err := fmt.Errorf("domain does not exist")
		_ = WriteResponse(w, http.StatusBadRequest, NewErrorPayload(404001, err.Error(), err))
		return
	}

	domainType := h.searchTypeDomain(domainName)

	if domainType == "unknown" {
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
					ClientIp: r.RemoteAddr},
				&result,
			)
		}(ctx)

		select {
		case <-ctx.Done():
			// If 1 second timeout reached, send a partial response and continue waiting for result
			duration := time.Since(start)
			partialPayload := NewResponsePayload(fmt.Sprintf("%s is defined as %s per %d milliseconds", requestPayload.Data, domainType, duration.Milliseconds()), domainType)
			_ = WriteResponse(w, http.StatusOK, partialPayload)
			h.log.WithFields(logrus.Fields{
				"rawData":  requestPayload.Data,
				"domain":   domainName,
				"type":     domainType,
				"duration": duration.Milliseconds(),
				"source":   "lookup",
			}).Info("successfully checked domain in lookup service via RPC")
		case err := <-errChan:
			if err == nil && validators.ValidateDomainTypes(result) == nil {
				domain := &models.Domain{
					Name:     domainName,
					Type:     result,
					Coverage: "equals",
				}
				err := h.domainRepo.Create(domain)
				if err != nil {
					h.log.Error(err)
				}
			} else {
				h.log.Error(err)
			}
		}
	}

	duration := time.Since(start)
	payload := NewResponsePayload(fmt.Sprintf("%s is defined as %s per %d milliseconds", requestPayload.Data, domainType, duration.Milliseconds()), domainType)
	_ = WriteResponse(w, http.StatusOK, payload)
	h.log.WithFields(logrus.Fields{
		"rawData":  requestPayload.Data,
		"domain":   domainName,
		"type":     domainType,
		"duration": duration.Milliseconds(),
		"source":   "database",
	}).Info("successfully checked domain in local database")
	return
}

func (h *BaseHandler) searchTypeDomain(domainName string) string {

	domainType := "unknown"

	chMatchEquals := make(chan string)
	chMatchContains := make(chan string)
	chMatchBegins := make(chan string)
	chMatchEnds := make(chan string)
	chQuit := make(chan bool)

	go func() {
		res, _ := h.domainRepo.MatchEquals(domainName)
		if res != nil {
			chMatchEquals <- res.Type
		}
		chQuit <- true
	}()

	go func() {
		res, _ := h.domainRepo.MatchContains(domainName)
		if res != nil {
			chMatchContains <- res.Type
		}
		chQuit <- true
	}()

	go func() {
		res, _ := h.domainRepo.MatchBegins(domainName)
		if res != nil {
			chMatchBegins <- res.Type
		}
		chQuit <- true
	}()

	go func() {
		res, _ := h.domainRepo.MatchEnds(domainName)
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
