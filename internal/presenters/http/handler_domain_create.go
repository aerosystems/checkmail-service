package HttpServer

import (
	CustomErrors "github.com/aerosystems/checkmail-service/internal/common/custom_errors"
	"github.com/labstack/echo/v4"
	"net/http"
)

type CreateDomainRequest struct {
	CreateDomainRequestBody
}

type CreateDomainRequestBody struct {
	Name     string `json:"name" validate:"fqdn,required" example:"gmail.com"`
	Type     string `json:"type" validate:"oneof=blacklist whitelist undefined, required" example:"whitelist"`
	Coverage string `json:"coverage" validate:"oneof=begins ends equals contains, required" example:"equals"`
}

// CreateDomain godoc
// @Summary create domain
// @Tags domains
// @Accept  json
// @Produce application/json
// @Param comment body CreateDomainRequestBody true "raw request body"
// @Security BearerAuth
// @Success 201 {object} Domain
// @Failure 400 {object} echo.HTTPError
// @Failure 401 {object} echo.HTTPError
// @Failure 403 {object} echo.HTTPError
// @Failure 409 {object} echo.HTTPError
// @Failure 422 {object} echo.HTTPError
// @Failure 500 {object} echo.HTTPError
// @Router /v1/domains [post]
func (dh DomainHandler) CreateDomain(c echo.Context) error {
	var requestPayload CreateDomainRequest
	if err := c.Bind(&requestPayload); err != nil {
		return CustomErrors.ErrInvalidRequestBody
	}
	domain, err := dh.domainUsecase.CreateDomain(requestPayload.Name, requestPayload.Type, requestPayload.Coverage)
	if err != nil {
		return err
	}
	return c.JSON(http.StatusCreated, ModelToDomain(domain))
}
