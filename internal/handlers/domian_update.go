package handlers

import (
	"errors"
	"github.com/aerosystems/checkmail-service/pkg/validators"
	"github.com/go-chi/chi/v5"
	"gorm.io/gorm"
	"net/http"
)

type UpdateDomainRequest struct {
	Name     string `json:"name" example:"gmail.com"`
	Type     string `json:"type" example:"whitelist"`
	Coverage string `json:"coverage" example:"equals"`
}

// DomainUpdate godoc
// @Summary update domain by Domain Name
// @Tags domains
// @Accept  json
// @Produce application/json
// @Param	domainName	path	string	true "Domain Name"
// @Param comment body models.Domain true "raw request body"
// @Security BearerAuth
// @Success 200 {object} Response{data=models.Domain}
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Failure 403 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Failure 422 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /v1/domains/{domainName} [patch]
func (h *BaseHandler) DomainUpdate(w http.ResponseWriter, r *http.Request) {
	domainName := chi.URLParam(r, "domainName")
	if domainName == "" {
		err := errors.New("path parameter domainName is empty")
		_ = WriteResponse(w, http.StatusBadRequest, NewErrorPayload(400202, err.Error(), err))
		return
	}

	var requestPayload UpdateDomainRequest

	if err := ReadRequest(w, r, &requestPayload); err != nil {
		_ = WriteResponse(w, http.StatusUnprocessableEntity, NewErrorPayload(422201, "could not read request body", err))
		return
	}

	if requestPayload == (UpdateDomainRequest{}) {
		err := errors.New("request payload could not be empty")
		_ = WriteResponse(w, http.StatusUnprocessableEntity, NewErrorPayload(422201, err.Error(), err))
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
		err := errors.New("name could not be changed")
		_ = WriteResponse(w, http.StatusUnprocessableEntity, NewErrorPayload(422205, err.Error(), err))
		return
	}

	if requestPayload.Type != "" {
		if err := validators.ValidateDomainTypes(requestPayload.Type); err != nil {
			_ = WriteResponse(w, http.StatusUnprocessableEntity, NewErrorPayload(422203, err.Error(), err))
			return
		}
		domain.Type = requestPayload.Type
	}

	if requestPayload.Coverage != "" {
		if err := validators.ValidateDomainCoverages(requestPayload.Coverage); err != nil {
			_ = WriteResponse(w, http.StatusUnprocessableEntity, NewErrorPayload(422204, err.Error(), err))
			return
		}
		domain.Coverage = requestPayload.Coverage
	}

	if err := h.domainRepo.Update(domain); err != nil {
		_ = WriteResponse(w, http.StatusInternalServerError, NewErrorPayload(500204, "could not update Domain by domainName", err))
		return
	}

	_ = WriteResponse(w, http.StatusOK, NewResponsePayload("domain successfully updated", domain))
	return
}
