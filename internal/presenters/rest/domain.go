package rest

import (
	"errors"
	"github.com/aerosystems/checkmail-service/internal/helpers"
	"github.com/aerosystems/checkmail-service/internal/models"
	CustomError "github.com/aerosystems/checkmail-service/pkg/custom_error"
	"github.com/aerosystems/checkmail-service/pkg/validators"
	"github.com/labstack/echo/v4"
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
func (h *BaseHandler) CreateDomain(c echo.Context) error {
	var requestPayload DomainRequest
	if err := c.Bind(&requestPayload); err != nil {
		return h.ErrorResponse(c, http.StatusUnprocessableEntity, "could not read request body", err)
	}
	if err := requestPayload.Validate(); err != nil {
		return h.ErrorResponse(c, http.StatusBadRequest, err.Message, err.Error())
	}
	root, _ := helpers.GetRootDomain(requestPayload.Name)
	rootDomain, _ := h.rootDomainRepo.FindByName(root)
	if rootDomain == nil {
		err := errors.New("domain does not exist")
		return h.ErrorResponse(c, http.StatusNotFound, err.Error(), err)
	}
	newDomain := models.Domain{
		Name:     requestPayload.Name,
		Type:     requestPayload.Type,
		Coverage: requestPayload.Coverage,
	}
	if err := h.domainRepo.Create(&newDomain); err != nil {
		if errors.Is(err, gorm.ErrDuplicatedKey) || strings.Contains(err.Error(), "duplicate key value violates unique constraint") {
			return h.ErrorResponse(c, http.StatusConflict, "domain already exists", err)
		}
		return h.ErrorResponse(c, http.StatusInternalServerError, "could not create new domain", err)
	}
	return h.SuccessResponse(c, http.StatusOK, "domain successfully created", newDomain)
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
func (h *BaseHandler) GetDomain(c echo.Context) error {
	domainName := c.Param("domainName")
	domain, err := h.domainRepo.FindByName(domainName)
	if err != nil {
		return h.ErrorResponse(c, http.StatusInternalServerError, "could not find domain", err)
	}
	if domain == nil {
		err := errors.New("domain not found")
		return h.ErrorResponse(c, http.StatusNotFound, err.Error(), err)
	}
	return h.SuccessResponse(c, http.StatusOK, "domain successfully found", domain)
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
func (h *BaseHandler) UpdateDomain(c echo.Context) error {
	var requestPayload DomainRequest
	if err := c.Bind(&requestPayload); err != nil {
		return h.ErrorResponse(c, http.StatusUnprocessableEntity, "could not read request body", err)
	}
	if err := requestPayload.Validate(); err != nil {
		return h.ErrorResponse(c, http.StatusBadRequest, err.Message, err.Error())
	}
	domainName := c.Param("domainName")
	domain, err := h.domainRepo.FindByName(domainName)
	if err != nil {
		return h.ErrorResponse(c, http.StatusInternalServerError, "could not find domain", err)
	}
	if domain == nil {
		err := errors.New("domain not found")
		return h.ErrorResponse(c, http.StatusNotFound, err.Error(), err)
	}
	if requestPayload.Name != "" {
		err := errors.New("name could not be changed")
		return h.ErrorResponse(c, http.StatusUnprocessableEntity, err.Error(), err)
	}
	if requestPayload.Type != "" {
		if err := validators.ValidateDomainTypes(requestPayload.Type); err != nil {
			return h.ErrorResponse(c, http.StatusUnprocessableEntity, err.Message, err.Error())
		}
		domain.Type = requestPayload.Type
	}
	if requestPayload.Coverage != "" {
		if err := validators.ValidateDomainCoverage(requestPayload.Coverage); err != nil {
			return h.ErrorResponse(c, http.StatusUnprocessableEntity, err.Message, err.Error())
		}
		domain.Coverage = requestPayload.Coverage
	}
	if err := h.domainRepo.Update(domain); err != nil {
		return h.ErrorResponse(c, http.StatusInternalServerError, "could not update domain", err)
	}
	return h.SuccessResponse(c, http.StatusOK, "domain successfully updated", domain)
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
func (h *BaseHandler) DeleteDomain(c echo.Context) error {
	domainName := c.Param("domainName")
	domain, err := h.domainRepo.FindByName(domainName)
	if err != nil {
		return h.ErrorResponse(c, http.StatusInternalServerError, "could not find domain", err)
	}
	if domain == nil {
		err := errors.New("domain not found")
		return h.ErrorResponse(c, http.StatusNotFound, err.Error(), err)
	}
	if err := h.domainRepo.Delete(domain); err != nil {
		return h.ErrorResponse(c, http.StatusInternalServerError, "could not delete domain", err)
	}
	return h.SuccessResponse(c, http.StatusOK, "domain successfully deleted", nil)
}
