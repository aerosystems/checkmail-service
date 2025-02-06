package HTTPServer

import (
	"github.com/labstack/echo/v4"
	"net/http"
)

type CreateFilterRequest struct {
	Name         string `json:"name" example:"gmail.com"`
	Type         string `json:"type" example:"whitelist"`
	Coverage     string `json:"coverage" example:"equals"`
	ProjectToken string `json:"projectToken" example:"38fa45ebb919g5d966122bf9g42a38ceb1e4f6eddf1da70ef00afbdc38197d8f"`
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
func (h Handler) CreateFilter(c echo.Context) error {
	return echo.NewHTTPError(http.StatusNotImplemented, "not implemented")
}
