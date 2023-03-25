package handlers

import (
	"context"
	"errors"
	"fmt"
	"github.com/aerosystems/checkmail-service/internal/helpers"
	"github.com/go-chi/chi/v5"
	"net/http"
	"net/mail"
	"strings"
	"time"
)

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
			_ = WriteResponse(w, http.StatusBadRequest, NewErrorPayload(400207, err.Error(), err))
			return
		}
		arr := strings.Split(email.Address, "@")
		domainName = arr[1]
	case 0:
		domainName = data
	default:
		err := errors.New("path param could not contain more then one \"@\" character")
		_ = WriteResponse(w, http.StatusBadRequest, NewErrorPayload(400208, err.Error(), err))
		return
	}

	// Validate Domain Name
	isValid := helpers.ValidateDomain(domainName)
	if !isValid {
		err := errors.New("domain does not valid")
		_ = WriteResponse(w, http.StatusBadRequest, NewErrorPayload(400210, err.Error(), err))
		return
	}

	domainType := "unknown"

	chMatchEquals := make(chan string)
	chMatchContains := make(chan string)
	chMatchBegins := make(chan string)
	chMatchEnds := make(chan string)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	//var wg sync.WaitGroup
	//wg.Add(4)

	go func(ctx context.Context, ch chan string) {
		//defer wg.Done()
		select {
		case <-ctx.Done():
			fmt.Println("finish equals")
			return
		default:
			time.Sleep(2 * time.Second)
			res, _ := h.domainRepo.MatchEquals(domainName)
			if res != nil {
				ch <- res.Type
				cancel()
			}
			return
		}
	}(ctx, chMatchEquals)

	go func(ctx context.Context, ch chan string) {
		//defer wg.Done()
		select {
		case <-ctx.Done():
			fmt.Println("finish contains")
			return
		default:
			time.Sleep(4 * time.Second)
			res, _ := h.domainRepo.MatchContains(domainName)
			if res != nil {
				ch <- res.Type
				cancel()
			}
			return
		}
	}(ctx, chMatchContains)

	go func(ctx context.Context, ch chan string) {
		//defer wg.Done()
		select {
		case <-ctx.Done():
			fmt.Println("finish begins")
			return
		default:
			time.Sleep(6 * time.Second)
			res, _ := h.domainRepo.MatchBegins(domainName)
			if res != nil {
				ch <- res.Type
				cancel()
			}
			return
		}
	}(ctx, chMatchBegins)

	go func(ctx context.Context, ch chan string) {
		//defer wg.Done()
		select {
		case <-ctx.Done():
			fmt.Println("finish ends")
			return
		default:
			time.Sleep(8 * time.Second)
			res, _ := h.domainRepo.MatchEnds(domainName)
			if res != nil {
				ch <- res.Type
				cancel()
			}
			return
		}
	}(ctx, chMatchEnds)

	go func(domainType string) {
		select {
		case domainType = <-chMatchEquals:
			fmt.Printf("find Equals first: %s", domainType)
		case domainType = <-chMatchContains:
			fmt.Printf("find Contains first: %s", domainType)
		case domainType = <-chMatchBegins:
			fmt.Printf("find Begins first: %s", domainType)
		case domainType = <-chMatchEnds:
			fmt.Printf("find Ends first: %s", domainType)
		}
	}(domainType)

	//wg.Wait()

	duration := time.Since(start)
	payload := NewResponsePayload(fmt.Sprintf("%s is defined as %s per %d milliseconds", data, domainType, duration.Milliseconds()), domainType)
	_ = WriteResponse(w, http.StatusOK, payload)
	return
}
