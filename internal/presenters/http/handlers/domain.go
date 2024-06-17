package handlers

import (
	"errors"
	CustomErrors "github.com/aerosystems/checkmail-service/internal/common/custom_errors"
	"github.com/labstack/echo/v4"
	"net/http"
)

type DomainHandler struct {
	*BaseHandler
	domainUsecase DomainUsecase
}

func NewDomainHandler(baseHandler *BaseHandler, domainUsecase DomainUsecase) *DomainHandler {
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
		return dh.ErrorResponse(c, CustomErrors.ErrReadRequestBody.HttpCode, CustomErrors.ErrReadRequestBody.Message, err)
	}
	domain, err := dh.domainUsecase.CreateDomain(requestPayload.Name, requestPayload.Type, requestPayload.Coverage)
	if err != nil {
		return dh.ErrorResponse(c, CustomErrors.ErrDomainInternalCreate.HttpCode, CustomErrors.ErrDomainInternalCreate.Message, err)
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
		var apiErr CustomErrors.ApiError
		if errors.As(err, &apiErr) {
			return dh.ErrorResponse(c, apiErr.HttpCode, apiErr.Message, err)
		}
		return dh.ErrorResponse(c, CustomErrors.ErrDomainInternalGet.HttpCode, CustomErrors.ErrDomainInternalGet.Message, err)
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
// @Success 200 {object} Response
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
		return dh.ErrorResponse(c, CustomErrors.ErrInvalidDomain.HttpCode, CustomErrors.ErrInvalidDomain.Message, err)
	}
	if err := dh.domainUsecase.UpdateDomain(domainName, requestPayload.Type, requestPayload.Coverage); err != nil {
		var apiErr CustomErrors.ApiError
		if errors.As(err, &apiErr) {
			return dh.ErrorResponse(c, apiErr.HttpCode, apiErr.Message, err)
		}
		return dh.ErrorResponse(c, CustomErrors.ErrDomainInternalUpdate.HttpCode, CustomErrors.ErrDomainInternalUpdate.Message, err)
	}
	return dh.SuccessResponse(c, http.StatusOK, "domain successfully updated", nil)
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
	if err := dh.domainUsecase.DeleteDomain(domainName); err != nil {
		var apiErr CustomErrors.ApiError
		if errors.As(err, &apiErr) {
			return dh.ErrorResponse(c, apiErr.HttpCode, apiErr.Message, err)
		}
		return dh.ErrorResponse(c, CustomErrors.ErrDomainInternalDelete.HttpCode, CustomErrors.ErrDomainInternalDelete.Message, err)
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
		var apiErr CustomErrors.ApiError
		if errors.As(err, &apiErr) {
			return dh.ErrorResponse(c, apiErr.HttpCode, apiErr.Message, err)
		}
		return dh.ErrorResponse(c, CustomErrors.ErrDomainInternalCount.HttpCode, CustomErrors.ErrDomainInternalCount.Message, err)
	}
	return dh.SuccessResponse(c, http.StatusOK, "domains successfully counted", count)
}
