package filter

import (
	"github.com/aerosystems/checkmail-service/internal/models"
	"github.com/aerosystems/checkmail-service/internal/validators"
	"github.com/labstack/echo/v4"
	"net/http"
)

type CreateFilterRequest struct {
	Name         string `json:"name" example:"gmail.com"`
	Type         string `json:"type" example:"whitelist"`
	Coverage     string `json:"coverage" example:"equals"`
	ProjectToken string `json:"projectToken" example:"38fa45ebb919g5d966122bf9g42a38ceb1e4f6eddf1da70ef00afbdc38197d8f"`
}

func (cr *CreateFilterRequest) Validate() *models.Error {
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
// @Success 201 {object} Response
// @Failure 400 {object} ErrorResponse
// @Failure 403 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Failure 409 {object} ErrorResponse
// @Failure 422 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /v1/filters [post]
func (fh Handler) CreateFilter(c echo.Context) error {
	var requestPayload CreateFilterRequest
	if err := c.Bind(&requestPayload); err != nil {
		return fh.ErrorResponse(c, http.StatusUnprocessableEntity, "could not read request body", err)
	}
	if err := requestPayload.Validate(); err != nil {
		return fh.ErrorResponse(c, http.StatusBadRequest, err.Message, err.Error())
	}
	filter, err := fh.filterUsecase.CreateFilter(requestPayload.Name, requestPayload.Type, requestPayload.Coverage, requestPayload.ProjectToken)
	if err != nil {
		return fh.ErrorResponse(c, http.StatusInternalServerError, "could not create filter", err)
	}
	return fh.SuccessResponse(c, http.StatusCreated, "filter successfully created", filter)
}
