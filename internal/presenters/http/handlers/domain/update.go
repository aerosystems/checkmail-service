package domain

import (
	CustomErrors "github.com/aerosystems/checkmail-service/internal/common/custom_errors"
	"github.com/labstack/echo/v4"
	"net/http"
)

type UpdateDomainRequest struct {
	UpdateDomainBody
	UpdateDomainQueryParam
}

type UpdateDomainBody struct {
	Type     string `json:"type" validate:"oneof=blacklist whitelist undefined, required" example:"whitelist"`
	Coverage string `json:"coverage" validate:"oneof=begins ends equals contains, required" example:"equals"`
}

type UpdateDomainQueryParam struct {
	Name string `json:"name" validate:"fqdn,required" example:"gmail.com"`
}

// UpdateDomain godoc
// @Summary update domain by Domain Name
// @Tags domains
// @Accept  json
// @Produce application/json
// @Param	domainName	path	string	true "Domain Name"
// @Param comment body UpdateDomainRequest true "raw request body"
// @Security BearerAuth
// @Success 200 {object} Domain
// @Failure 400 {object} echo.HTTPError
// @Failure 401 {object} echo.HTTPError
// @Failure 403 {object} echo.HTTPError
// @Failure 404 {object} echo.HTTPError
// @Failure 422 {object} echo.HTTPError
// @Failure 500 {object} echo.HTTPError
// @Router /v1/domains/{domainName} [patch]
func (dh Handler) UpdateDomain(c echo.Context) error {
	var requestPayload CreateDomainRequestBody
	if err := c.Bind(&requestPayload); err != nil {
		return CustomErrors.ErrInvalidRequestBody
	}
	domain, err := dh.domainUsecase.UpdateDomain(requestPayload.Name, requestPayload.Type, requestPayload.Coverage)
	if err != nil {
		return err
	}
	return c.JSON(http.StatusOK, ModelToDomain(domain))
}
