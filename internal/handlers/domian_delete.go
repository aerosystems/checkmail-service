package handlers

import (
	"errors"
	"github.com/go-chi/chi/v5"
	"gorm.io/gorm"
	"net/http"
)

// DomainDelete godoc
// @Summary delete domain by Domain Name
// @Tags domains
// @Accept  json
// @Produce application/json
// @Param	domainName	path	string	true "Domain Name"
// @Security BearerAuth
// @Success 200 {object} Response
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Failure 403 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /v1/domains/{domainName} [delete]
func (h *BaseHandler) DomainDelete(w http.ResponseWriter, r *http.Request) {
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

	if err := h.domainRepo.Delete(domain); err != nil {
		_ = WriteResponse(w, http.StatusInternalServerError, NewErrorPayload(500202, "could not delete Domain", err))
		return
	}

	_ = WriteResponse(w, http.StatusOK, NewResponsePayload("domain successfully deleted", nil))
	return
}
