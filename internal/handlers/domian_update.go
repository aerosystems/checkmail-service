package handlers

import (
	"errors"
	"github.com/aerosystems/checkmail-service/internal/models"
	"github.com/aerosystems/checkmail-service/pkg/validators"
	"github.com/go-chi/chi/v5"
	"gorm.io/gorm"
	"net/http"
)

// DomainUpdate godoc
// @Summary update domain by Domain Name
// @Tags domains
// @Accept  json
// @Produce application/json
// @Param	domainName	path	string	true "Domain Name"
// @Param comment body models.Domain true "raw request body"
// @Param Authorization header string true "should contain Access Token, with the Bearer started"
// @Success 200 {object} Response{data=models.Domain}
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /domains/{domainName} [patch]
func (h *BaseHandler) DomainUpdate(w http.ResponseWriter, r *http.Request) {
	domainName := chi.URLParam(r, "domainName")
	if domainName == "" {
		err := errors.New("path parameter domainName is empty")
		_ = WriteResponse(w, http.StatusBadRequest, NewErrorPayload(400202, err.Error(), err))
		return
	}

	var requestPayload models.Domain

	if err := ReadRequest(w, r, &requestPayload); err != nil {
		_ = WriteResponse(w, http.StatusBadRequest, NewErrorPayload(400201, "request payload is incorrect", err))
		return
	}

	if requestPayload == (models.Domain{}) {
		err := errors.New("request payload could not be empty")
		_ = WriteResponse(w, http.StatusInternalServerError, NewErrorPayload(500201, err.Error(), err))
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

	if requestPayload.Name != "" {
		err := errors.New("claim Name could not be changed")
		_ = WriteResponse(w, http.StatusBadRequest, NewErrorPayload(400205, err.Error(), err))
		return
	}

	if requestPayload.Type != "" {
		if err := validators.ValidateDomainTypes(requestPayload.Type); err != nil {
			_ = WriteResponse(w, http.StatusBadRequest, NewErrorPayload(400203, err.Error(), err))
			return
		}
		domain.Type = requestPayload.Type
	}

	if requestPayload.Coverage != "" {
		if err := validators.ValidateDomainCoverages(requestPayload.Coverage); err != nil {
			_ = WriteResponse(w, http.StatusBadRequest, NewErrorPayload(400204, err.Error(), err))
			return
		}
		domain.Coverage = requestPayload.Coverage
	}

	if err := h.domainRepo.Update(domain); err != nil {
		_ = WriteResponse(w, http.StatusInternalServerError, NewErrorPayload(500204, "could not update Domain by domainName", err))
		return
	}

	payload := NewResponsePayload("domain successfully updated", domain)
	_ = WriteResponse(w, http.StatusNotImplemented, payload)
	return
}
