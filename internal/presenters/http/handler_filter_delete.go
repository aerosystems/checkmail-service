package HTTPServer

import (
	"github.com/labstack/echo/v4"
	"net/http"
)

// DeleteFilter godoc
// @Summary Delete Filter
// @Description Delete Filter for ProjectRPCPayload by projectId. Roles allowed: business, staff
// @Tags Filter
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param id query int true "filter id"
// @Success 204
// @Failure 400 {object} echo.HTTPError
// @Failure 403 {object} echo.HTTPError
// @Failure 404 {object} echo.HTTPError
// @Failure 409 {object} echo.HTTPError
// @Failure 422 {object} echo.HTTPError
// @Failure 500 {object} echo.HTTPError
// @Router /v1/filters/{domain_name} [delete]
func (fh FilterHandler) DeleteFilter(c echo.Context) error {
	return echo.NewHTTPError(http.StatusNotImplemented, "not implemented")
}
