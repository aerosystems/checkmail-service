package handlers

import (
	"errors"
	"fmt"
	"github.com/aerosystems/checkmail-service/internal/helpers"
	"github.com/go-chi/chi/v5"
	"net/http"
	"net/mail"
	"strings"
)

func (h *BaseHandler) Data(w http.ResponseWriter, r *http.Request) {
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
	fmt.Println(domains)

	payload := NewResponsePayload("method not implemented", nil)
	_ = WriteResponse(w, http.StatusNotImplemented, payload)
	return
}
