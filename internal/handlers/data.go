package handlers

import (
	"errors"
	"fmt"
	"github.com/go-chi/chi/v5"
	"golang.org/x/net/idna"
	"gorm.io/gorm"
	"net/http"
	"net/mail"
	"strings"
)

func (h *BaseHandler) Data(w http.ResponseWriter, r *http.Request) {
	data := chi.URLParam(r, "data")
	data = strings.ToLower(data)

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

	he, err := idna.Lookup.ToASCII(domainName)
	fmt.Println(he, err)

	domain, err := h.domainRepo.FindByName(domainName)
	if err != nil && err != gorm.ErrRecordNotFound {
		_ = WriteResponse(w, http.StatusInternalServerError, NewErrorPayload(500201, "could not find Domain by domainName", err))
		return
	}
	if domain == nil {
		// TODO send request for lookup MX records
		payload := NewResponsePayload("method not implemented", nil)
		_ = WriteResponse(w, http.StatusNotImplemented, payload)
		return
	}

	payload := NewResponsePayload("method not implemented", nil)
	_ = WriteResponse(w, http.StatusNotImplemented, payload)
	return
}
