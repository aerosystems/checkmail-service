package HTTPServer

import (
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
// @Success 200 {object} Domain
// @Failure 400 {object} echo.HTTPError
// @Failure 401 {object} echo.HTTPError
// @Failure 403 {object} echo.HTTPError
// @Failure 404 {object} echo.HTTPError
// @Failure 500 {object} echo.HTTPError
// @Router /v1/domains/{domainName} [get]
func (dh DomainHandler) GetDomain(c echo.Context) error {
	var requestPayload GetDomainRequest
	if err := c.Bind(&requestPayload); err != nil {
		return CustomErrors.ErrInvalidRequestBody
	}
	domain, err := dh.domainUsecase.GetDomainByName(requestPayload.Name)
	if err != nil {
		return err
	}
	return c.JSON(http.StatusOK, ModelToDomain(domain))
}
