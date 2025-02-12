package HTTPServer

import (
	"github.com/aerosystems/checkmail-service/internal/entities"
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
// @Router /v1/domains/{domain_name} [delete]
func (h Handler) DeleteDomain(c echo.Context) error {
	var requestPayload DeleteDomainRequest
	if err := c.Bind(&requestPayload); err != nil {
		return entities.ErrInvalidRequestBody
	}
	if err := h.domainUsecase.DeleteDomain(c.Request().Context(), requestPayload.Name); err != nil {
		return err
	}
	return c.JSON(http.StatusNoContent, nil)
}
