package handlers

import (
	"context"
	"errors"
	"fmt"
	"github.com/aerosystems/checkmail-service/internal/helpers"
	"github.com/aerosystems/checkmail-service/internal/models"
	"github.com/go-chi/chi/v5"
	"net/http"
	"net/mail"
	"strings"
	"sync"
	"time"
)

type DataResponse struct {
	Type                string `json:"type"`
	ExecutionTimeMillis int64
}

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

	// Check Domain Name
	domains, _ := h.domainRepo.FindAll()
	typeDomain := "unknown"

	result := make(chan string)
	var wg sync.WaitGroup
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	for _, item := range *domains {
		wg.Add(1)
		item := item
		go func() {
			search(ctx, &item, domainName, result)
			wg.Done()
		}()
	}

	go func() {
		typeDomain = <-result
		cancel()
	}()

	wg.Wait()
	// Отримання результатів роботи горутин
	fmt.Printf("Received result %s\n", typeDomain)

	duration := time.Since(start)
	payload := NewResponsePayload("domain successfully detected", DataResponse{
		Type:                typeDomain,
		ExecutionTimeMillis: duration.Milliseconds(),
	})
	_ = WriteResponse(w, http.StatusOK, payload)
	return
}

func search(ctx context.Context, item *models.Domain, domainName string, result chan string) {
	for {
		select {
		case <-ctx.Done():
			return
		default:
			fmt.Printf("%s : %s : %t", item.Name, item.Coverage, item.Match(domainName))
			if item.Match(domainName) {
				result <- item.Type
				return
			}
			continue
		}
	}
}
