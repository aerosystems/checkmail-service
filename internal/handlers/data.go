package handlers

import (
	"errors"
	"fmt"
	"github.com/aerosystems/checkmail-service/pkg/validators"
	"github.com/go-chi/chi/v5"
	"net/http"
	"net/mail"
	"strings"
	"sync"
	"time"
)

// Data godoc
// @Summary get information about domain name or email address
// @Tags data
// @Accept  json
// @Produce application/json
// @Param	data	path	string	true "Domain Name or Email Address"
// @Security X-API-KEY
// @Success 200 {object} Response{data=string}
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Failure 422 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /v1/data/{data} [get]
func (h *BaseHandler) Data(w http.ResponseWriter, r *http.Request) {
	start := time.Now()
	data := chi.URLParam(r, "data")
	data = strings.ToLower(data)

	// Get Domain Name
	var domainName string
	switch strings.Count(data, "@") {
	case 1:
		email, err := mail.ParseAddress(data)
		if err != nil {
			_ = WriteResponse(w, http.StatusUnprocessableEntity, NewErrorPayload(422207, err.Error(), err))
			return
		}
		arr := strings.Split(email.Address, "@")
		domainName = arr[1]
	case 0:
		domainName = data
	default:
		err := errors.New("path param could not contain more then one \"@\" character")
		_ = WriteResponse(w, http.StatusUnprocessableEntity, NewErrorPayload(422208, err.Error(), err))
		return
	}

	// Validate Domain Name
	isValid := validators.ValidateDomain(domainName)
	if !isValid {
		err := errors.New("domain does not valid")
		_ = WriteResponse(w, http.StatusUnprocessableEntity, NewErrorPayload(422210, err.Error(), err))
		return
	}

	// Check Root Domain Name
	arrDomain := strings.Split(domainName, ".")
	root := arrDomain[len(arrDomain)-1]
	rootDomain, _ := h.rootDomainRepo.FindByName(root)
	if rootDomain == nil {
		err := fmt.Errorf("domain '%s' does not exist, because '%s' is not root domain", domainName, root)
		_ = WriteResponse(w, http.StatusBadRequest, NewErrorPayload(400211, err.Error(), err))
		return
	}

	domainType := h.SearchTypeDomain(domainName)

	duration := time.Since(start)
	payload := NewResponsePayload(fmt.Sprintf("%s is defined as %s per %d milliseconds", data, domainType, duration.Milliseconds()), domainType)
	_ = WriteResponse(w, http.StatusOK, payload)
	return
}

func (h *BaseHandler) SearchTypeDomain(domainName string) string {

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
