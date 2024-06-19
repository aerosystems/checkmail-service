package domain

import (
	"errors"
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
// @Success 200 {object} Response
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Failure 403 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Failure 422 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /v1/domains/{domainName} [patch]
func (dh Handler) UpdateDomain(c echo.Context) error {
	var requestPayload CreateDomainRequestBody
	if err := c.Bind(&requestPayload); err != nil {
		return dh.ErrorResponse(c, http.StatusUnprocessableEntity, "could not read request body", err)
	}
	if err := dh.domainUsecase.UpdateDomain(requestPayload.Name, requestPayload.Type, requestPayload.Coverage); err != nil {
		var apiErr CustomErrors.ApiError
		if errors.As(err, &apiErr) {
			return dh.ErrorResponse(c, apiErr.HttpCode, apiErr.Message, err)
		}
		return dh.ErrorResponse(c, CustomErrors.ErrDomainInternalUpdate.HttpCode, CustomErrors.ErrDomainInternalUpdate.Message, err)
	}
	return dh.SuccessResponse(c, http.StatusOK, "domain successfully updated", nil)
}
