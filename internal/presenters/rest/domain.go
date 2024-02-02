package rest

import (
	"errors"
	"github.com/labstack/echo/v4"
	"net/http"
)

type DomainHandler struct {
	BaseHandler
	domainUsecase DomainUsecase
}

func NewDomainHandler(baseHandler BaseHandler, domainUsecase DomainUsecase) *DomainHandler {
	return &DomainHandler{
		BaseHandler:   baseHandler,
		domainUsecase: domainUsecase,
	}
}

type CreateDomainRequest struct {
	Name     string `json:"name" validate:"fqdn, required" example:"gmail.com"`
	Type     string `json:"type" validate:"oneof=blacklist whitelist undefined, required" example:"whitelist"`
	Coverage string `json:"coverage" validate:"oneof=begins ends equals contains, required" example:"equals"`
}

type UpdateDomainRequest struct {
	Type     string `json:"type" validate:"oneof=blacklist whitelist undefined, required" example:"whitelist"`
	Coverage string `json:"coverage" validate:"oneof=begins ends equals contains, required" example:"equals"`
}

// CreateDomain godoc
// @Summary create domain
// @Tags domains
// @Accept  json
// @Produce application/json
// @Param comment body CreateDomainRequest true "raw request body"
// @Security BearerAuth
// @Success 200 {object} Response{data=models.Domain}
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Failure 403 {object} ErrorResponse
// @Failure 409 {object} ErrorResponse
// @Failure 422 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /v1/domains [post]
func (dh DomainHandler) CreateDomain(c echo.Context) error {
	var requestPayload CreateDomainRequest
	if err := c.Bind(&requestPayload); err != nil {
		return dh.ErrorResponse(c, http.StatusUnprocessableEntity, "could not read request body", err)
	}
	domain, err := dh.domainUsecase.CreateDomain(requestPayload.Name, requestPayload.Type, requestPayload.Coverage)
	if err != nil {
		return dh.ErrorResponse(c, http.StatusInternalServerError, "could not create domain", err)
	}
	return dh.SuccessResponse(c, http.StatusOK, "domain successfully created", domain)
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
func (dh DomainHandler) GetDomain(c echo.Context) error {
	domainName := c.Param("domainName")
	if err := dh.validator.Var(domainName, "required,fqdn"); err != nil {
		return dh.ErrorResponse(c, http.StatusBadRequest, "invalid domain name", err)
	}
	domain, err := dh.domainUsecase.GetDomainByName(domainName)
	if err != nil {
		return dh.ErrorResponse(c, http.StatusInternalServerError, "could not find domain", err)
	}
	if domain == nil {
		err := errors.New("domain not found")
		return dh.ErrorResponse(c, http.StatusNotFound, err.Error(), err)
	}
	return dh.SuccessResponse(c, http.StatusOK, "domain successfully found", domain)
}

// UpdateDomain godoc
// @Summary update domain by Domain Name
// @Tags domains
// @Accept  json
// @Produce application/json
// @Param	domainName	path	string	true "Domain Name"
// @Param comment body UpdateDomainRequest true "raw request body"
// @Security BearerAuth
// @Success 200 {object} Response{data=models.Domain}
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Failure 403 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Failure 422 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /v1/domains/{domainName} [patch]
func (dh DomainHandler) UpdateDomain(c echo.Context) error {
	var requestPayload CreateDomainRequest
	if err := c.Bind(&requestPayload); err != nil {
		return dh.ErrorResponse(c, http.StatusUnprocessableEntity, "could not read request body", err)
	}
	domainName := c.Param("domainName")
	if err := dh.validator.Var(domainName, "required,fqdn"); err != nil {
		return dh.ErrorResponse(c, http.StatusBadRequest, "invalid domain name", err)
	}
	domain, err := dh.domainUsecase.GetDomainByName(domainName)
	if err != nil {
		return dh.ErrorResponse(c, http.StatusInternalServerError, "could not find domain", err)
	}
	if domain == nil {
		err := errors.New("domain not found")
		return dh.ErrorResponse(c, http.StatusNotFound, err.Error(), err)
	}
	if err := dh.domainUsecase.UpdateDomain(domain, requestPayload.Type, requestPayload.Coverage); err != nil {
		return dh.ErrorResponse(c, http.StatusInternalServerError, "could not find domain", err)
	}
	return dh.SuccessResponse(c, http.StatusOK, "domain successfully updated", domain)
}

// DeleteDomain godoc
// @Summary delete domain by Domain Name
// @Tags domains
// @Accept  json
// @Produce application/json
// @Param	domainName	path	string	true "Domain Name"
// @Security BearerAuth
// @Success 204 {object} Response
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Failure 403 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /v1/domains/{domainName} [delete]
func (dh DomainHandler) DeleteDomain(c echo.Context) error {
	domainName := c.Param("domainName")
	if err := dh.validator.Var(domainName, "required,fqdn"); err != nil {
		return dh.ErrorResponse(c, http.StatusBadRequest, "invalid domain name", err)
	}
	domain, err := dh.domainUsecase.GetDomainByName(domainName)
	if err != nil {
		return dh.ErrorResponse(c, http.StatusInternalServerError, "could not find domain", err)
	}
	if domain == nil {
		err := errors.New("domain not found")
		return dh.ErrorResponse(c, http.StatusNotFound, err.Error(), err)
	}
	if err := dh.domainUsecase.DeleteDomain(domain); err != nil {
		return dh.ErrorResponse(c, http.StatusInternalServerError, "could not delete domain", err)
	}
	return dh.SuccessResponse(c, http.StatusNoContent, "domain successfully deleted", nil)
}

// Count godoc
// @Summary count Domains
// @Tags domains
// @Accept  json
// @Produce application/json
// @Success 200 {object} Response
// @Failure 500 {object} ErrorResponse
// @Router /v1/domains/count [post]
func (dh DomainHandler) Count(c echo.Context) error {
	count, err := dh.domainUsecase.CountDomains()
	if err != nil {
		return dh.ErrorResponse(c, http.StatusInternalServerError, "could not count Domains", err)
	}
	return dh.SuccessResponse(c, http.StatusOK, "domains successfully counted", count)
}
