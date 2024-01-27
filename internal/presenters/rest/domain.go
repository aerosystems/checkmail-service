package rest

import (
	"errors"
	"github.com/aerosystems/checkmail-service/internal/helpers"
	"github.com/aerosystems/checkmail-service/internal/models"
	CustomError "github.com/aerosystems/checkmail-service/pkg/custom_error"
	"github.com/aerosystems/checkmail-service/pkg/validators"
	"github.com/go-chi/chi/v5"
	"gorm.io/gorm"
	"net/http"
	"strings"
)

type DomainRequest struct {
	Name     string `json:"name" example:"gmail.com"`
	Type     string `json:"type" example:"whitelist"`
	Coverage string `json:"coverage" example:"equals"`
}

func (r *DomainRequest) Validate() *CustomError.Error {
	if err := validators.ValidateDomainTypes(r.Type); err != nil {
		return err
	}
	if err := validators.ValidateDomainCoverage(r.Coverage); err != nil {
		return err
	}
	if err := validators.ValidateDomainName(r.Name); err != nil {
		return err
	}
	return nil
}

// CreateDomain godoc
// @Summary create domain
// @Tags domains
// @Accept  json
// @Produce application/json
// @Param comment body DomainRequest true "raw request body"
// @Security BearerAuth
// @Success 200 {object} Response{data=models.Domain}
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Failure 403 {object} ErrorResponse
// @Failure 409 {object} ErrorResponse
// @Failure 422 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /v1/domains [post]
func (h *BaseHandler) CreateDomain(w http.ResponseWriter, r *http.Request) {
	var requestPayload DomainRequest
	if err := ReadRequest(w, r, &requestPayload); err != nil {
		_ = WriteResponse(w, http.StatusUnprocessableEntity, NewErrorPayload(422201, "could not read request body", err))
		return
	}
	if err := requestPayload.Validate(); err != nil {
		_ = WriteResponse(w, http.StatusBadRequest, NewErrorPayload(err.Code, err.Message, err.Error()))
		return
	}
	root, _ := helpers.GetRootDomain(requestPayload.Name)
	rootDomain, _ := h.rootDomainRepo.FindByName(root)
	if rootDomain == nil {
		err := errors.New("domain does not exist")
		_ = WriteResponse(w, http.StatusBadRequest, NewErrorPayload(400205, err.Error(), err))
		return
	}
	newDomain := models.Domain{
		Name:     requestPayload.Name,
		Type:     requestPayload.Type,
		Coverage: requestPayload.Coverage,
	}
	if err := h.domainRepo.Create(&newDomain); err != nil {
		if errors.Is(err, gorm.ErrDuplicatedKey) || strings.Contains(err.Error(), "duplicate key value violates unique constraint") {
			_ = WriteResponse(w, http.StatusConflict, NewErrorPayload(409202, "domain already exists", err))
			return
		}
		_ = WriteResponse(w, http.StatusInternalServerError, NewErrorPayload(500203, "could not create new domain", err))
		return
	}
	_ = WriteResponse(w, http.StatusOK, NewResponsePayload("domain successfully created", newDomain))
	return
}

// GetDomain godoc
// @Summary get domain by Domain Name
// @Tags domains
// @Accept  json
// @Produce application/json
// @Param	domainName	path	string	true "Domain Name"
// @Security BearerAuth
// @Success 200 {object} Response{data=models.Domain}
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Failure 403 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /v1/domains/{domainName} [get]
func (h *BaseHandler) GetDomain(w http.ResponseWriter, r *http.Request) {
	domainName := chi.URLParam(r, "domainName")
	domain, err := h.domainRepo.FindByName(domainName)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		_ = WriteResponse(w, http.StatusInternalServerError, NewErrorPayload(500201, "could not find domain", err))
		return
	}
	if domain == nil {
		err := errors.New("domain not found")
		_ = WriteResponse(w, http.StatusNotFound, NewErrorPayload(404201, err.Error(), err))
		return
	}
	_ = WriteResponse(w, http.StatusOK, NewResponsePayload("domain successfully found", domain))
	return
}

// UpdateDomain godoc
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
func (h *BaseHandler) UpdateDomain(w http.ResponseWriter, r *http.Request) {
	var requestPayload DomainRequest
	if err := ReadRequest(w, r, &requestPayload); err != nil {
		_ = WriteResponse(w, http.StatusUnprocessableEntity, NewErrorPayload(422201, "could not read request body", err))
		return
	}
	if err := requestPayload.Validate(); err != nil {
		_ = WriteResponse(w, http.StatusBadRequest, NewErrorPayload(err.Code, err.Message, err.Error()))
		return
	}
	domainName := chi.URLParam(r, "domainName")
	domain, err := h.domainRepo.FindByName(domainName)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		_ = WriteResponse(w, http.StatusInternalServerError, NewErrorPayload(500201, "could not find domain", err))
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
			_ = WriteResponse(w, http.StatusUnprocessableEntity, NewErrorPayload(err.Code, err.Message, err.Error()))
			return
		}
		domain.Type = requestPayload.Type
	}
	if requestPayload.Coverage != "" {
		if err := validators.ValidateDomainCoverage(requestPayload.Coverage); err != nil {
			_ = WriteResponse(w, http.StatusUnprocessableEntity, NewErrorPayload(err.Code, err.Message, err.Error()))
			return
		}
		domain.Coverage = requestPayload.Coverage
	}
	if err := h.domainRepo.Update(domain); err != nil {
		_ = WriteResponse(w, http.StatusInternalServerError, NewErrorPayload(500204, "could not update domain", err))
		return
	}
	_ = WriteResponse(w, http.StatusOK, NewResponsePayload("domain successfully updated", domain))
	return
}

// DeleteDomain godoc
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
func (h *BaseHandler) DeleteDomain(w http.ResponseWriter, r *http.Request) {
	domainName := chi.URLParam(r, "domainName")
	domain, err := h.domainRepo.FindByName(domainName)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
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
