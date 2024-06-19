package domain

import (
	"errors"
	CustomErrors "github.com/aerosystems/checkmail-service/internal/common/custom_errors"
	"github.com/labstack/echo/v4"
	"net/http"
)

type GetDomainRequest struct {
	UpdateDomainQueryParam
}

type GetDomainQueryParam struct {
	Name string `json:"name" validate:"fqdn,required" example:"gmail.com"`
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
func (dh Handler) GetDomain(c echo.Context) error {
	var requestPayload GetDomainRequest
	if err := c.Bind(&requestPayload); err != nil {
		return dh.ErrorResponse(c, CustomErrors.ErrInvalidDomain.HttpCode, CustomErrors.ErrInvalidDomain.Message, err)
	}
	domain, err := dh.domainUsecase.GetDomainByName(requestPayload.Name)
	if err != nil {
		var apiErr CustomErrors.ApiError
		if errors.As(err, &apiErr) {
			return dh.ErrorResponse(c, apiErr.HttpCode, apiErr.Message, err)
		}
		return dh.ErrorResponse(c, CustomErrors.ErrDomainInternalGet.HttpCode, CustomErrors.ErrDomainInternalGet.Message, err)
	}
	return dh.SuccessResponse(c, http.StatusOK, "domain successfully found", ModelToDomain(domain))
}
