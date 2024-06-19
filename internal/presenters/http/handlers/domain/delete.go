package domain

import (
	"errors"
	CustomErrors "github.com/aerosystems/checkmail-service/internal/common/custom_errors"
	"github.com/labstack/echo/v4"
	"net/http"
)

type DeleteDomainRequest struct {
	UpdateDomainQueryParam
}

type DeleteDomainQueryParam struct {
	Name string `json:"name" validate:"fqdn,required" example:"gmail.com"`
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
func (dh Handler) DeleteDomain(c echo.Context) error {
	var requestPayload DeleteDomainRequest
	if err := c.Bind(&requestPayload); err != nil {
		return dh.ErrorResponse(c, CustomErrors.ErrInvalidDomain.HttpCode, CustomErrors.ErrInvalidDomain.Message, err)
	}
	if err := dh.domainUsecase.DeleteDomain(requestPayload.Name); err != nil {
		var apiErr CustomErrors.ApiError
		if errors.As(err, &apiErr) {
			return dh.ErrorResponse(c, apiErr.HttpCode, apiErr.Message, err)
		}
		return dh.ErrorResponse(c, CustomErrors.ErrDomainInternalDelete.HttpCode, CustomErrors.ErrDomainInternalDelete.Message, err)
	}
	return dh.SuccessResponse(c, http.StatusNoContent, "domain successfully deleted", nil)
}
