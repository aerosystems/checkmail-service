package domain

import (
	CustomErrors "github.com/aerosystems/checkmail-service/internal/common/custom_errors"
	"github.com/labstack/echo/v4"
	"net/http"
)

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
// @Success 200 {object} Response{data=models.Domain}
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Failure 403 {object} ErrorResponse
// @Failure 409 {object} ErrorResponse
// @Failure 422 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /v1/domains [post]
func (dh Handler) CreateDomain(c echo.Context) error {
	var requestPayload CreateDomainRequestBody
	if err := c.Bind(&requestPayload); err != nil {
		return dh.ErrorResponse(c, CustomErrors.ErrReadRequestBody.HttpCode, CustomErrors.ErrReadRequestBody.Message, err)
	}
	domain, err := dh.domainUsecase.CreateDomain(requestPayload.Name, requestPayload.Type, requestPayload.Coverage)
	if err != nil {
		return dh.ErrorResponse(c, CustomErrors.ErrDomainInternalCreate.HttpCode, CustomErrors.ErrDomainInternalCreate.Message, err)
	}
	return dh.SuccessResponse(c, http.StatusOK, "domain successfully created", domain)
}
