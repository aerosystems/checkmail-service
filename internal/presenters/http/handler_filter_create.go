package HttpServer

import (
	CustomErrors "github.com/aerosystems/checkmail-service/internal/common/custom_errors"
	"github.com/aerosystems/checkmail-service/internal/common/validators"
	"github.com/labstack/echo/v4"
	"net/http"
)

type CreateFilterRequest struct {
	Name         string `json:"name" example:"gmail.com"`
	Type         string `json:"type" example:"whitelist"`
	Coverage     string `json:"coverage" example:"equals"`
	ProjectToken string `json:"projectToken" example:"38fa45ebb919g5d966122bf9g42a38ceb1e4f6eddf1da70ef00afbdc38197d8f"`
}

func (cr *CreateFilterRequest) Validate() error {
	if err := validators.ValidateDomainTypes(cr.Type); err != nil {
		return err
	}
	if err := validators.ValidateDomainCoverage(cr.Coverage); err != nil {
		return err
	}
	if err := validators.ValidateDomainName(cr.Name); err != nil {
		return err
	}
	return nil
}

// CreateFilter godoc
// @Summary Create Filter
// @Description Create Filter for ProjectRPCPayload. Roles allowed: business, staff
// @Tags Filter
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param filter body CreateFilterRequest true "raw request body"
// @Success 201 {object} Filter
// @Failure 400 {object} echo.HTTPError
// @Failure 403 {object} echo.HTTPError
// @Failure 404 {object} echo.HTTPError
// @Failure 409 {object} echo.HTTPError
// @Failure 422 {object} echo.HTTPError
// @Failure 500 {object} echo.HTTPError
// @Router /v1/filters [post]
func (fh FilterHandler) CreateFilter(c echo.Context) error {
	var requestPayload CreateFilterRequest
	if err := c.Bind(&requestPayload); err != nil {
		return CustomErrors.ErrInvalidRequestBody
	}
	if err := requestPayload.Validate(); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	filter, err := fh.filterUsecase.CreateFilter(requestPayload.Name, requestPayload.Type, requestPayload.Coverage, requestPayload.ProjectToken)
	if err != nil {
		return err
	}
	return c.JSON(http.StatusCreated, ModelToFilter(filter))
}
