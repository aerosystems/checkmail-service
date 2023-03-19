package handlers

import (
	"errors"
	"github.com/go-chi/chi/v5"
	"gorm.io/gorm"
	"net/http"
)

func (h *BaseHandler) DomainRead(w http.ResponseWriter, r *http.Request) {
	domainName := chi.URLParam(r, "domainName")
	if domainName == "" {
		err := errors.New("path parameter domainName is empty")
		_ = WriteResponse(w, http.StatusBadRequest, NewErrorPayload(400202, err.Error(), err))
		return
	}

	domain, err := h.domainRepo.FindByName(domainName)
	if err != nil && err != gorm.ErrRecordNotFound {
		_ = WriteResponse(w, http.StatusInternalServerError, NewErrorPayload(500201, "could not find Domain by domainName", err))
		return
	}

	if domain == nil {
		err := errors.New("domain not found")
		_ = WriteResponse(w, http.StatusNotFound, NewErrorPayload(404201, err.Error(), err))
		return
	}

	payload := NewResponsePayload("domain successfully found", domain)
	_ = WriteResponse(w, http.StatusOK, payload)
	return
}
