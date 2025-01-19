package HttpServer

import (
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
// @Success 204
// @Failure 400 {object} echo.HTTPError
// @Failure 401 {object} echo.HTTPError
// @Failure 403 {object} echo.HTTPError
// @Failure 404 {object} echo.HTTPError
// @Failure 500 {object} echo.HTTPError
// @Router /v1/domains/{domainName} [delete]
func (dh DomainHandler) DeleteDomain(c echo.Context) error {
	var requestPayload DeleteDomainRequest
	if err := c.Bind(&requestPayload); err != nil {
		return CustomErrors.ErrInvalidRequestBody
	}
	if err := dh.domainUsecase.DeleteDomain(requestPayload.Name); err != nil {
		return err
	}
	return c.JSON(http.StatusNoContent, nil)
}
