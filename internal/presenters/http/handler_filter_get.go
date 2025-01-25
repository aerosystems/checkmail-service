package HTTPServer

import (
	"github.com/labstack/echo/v4"
	"net/http"
)

// GetFilterList godoc
// @Summary Get Filter List
// @Description Get Filter List for all user Projects. Roles allowed: business, staff
// @Tags Filter
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param userId query int false "user id"
// @Param projectToken query string false "project token"
// @Success 200 {object} Filter
// @Failure 400 {object} echo.HTTPError
// @Failure 403 {object} echo.HTTPError
// @Failure 404 {object} echo.HTTPError
// @Failure 409 {object} echo.HTTPError
// @Failure 422 {object} echo.HTTPError
// @Failure 500 {object} echo.HTTPError
// @Router /v1/filters [get]
func (fh FilterHandler) GetFilterList(c echo.Context) error {
	return echo.NewHTTPError(http.StatusNotImplemented, "not implemented")
}
