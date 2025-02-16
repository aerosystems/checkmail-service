package HTTPServer

import (
	"github.com/labstack/echo/v4"
	"net/http"
)

// CheckData godoc
// @Summary Get DomainType about domain/email address
// @Description Get DomainType about domain/email address
// @Tags data
// @Accept  json
// @Produce application/json
// @Security BearerAuth
// @Success 200 {object} models.InspectResult
// @Failure 401 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /v1/domain/check [get]
func (h Handler) CheckData(c echo.Context) error {
	res, err := h.inspectUsecase.InspectDataDeprecated(
		c.Request().Context(),
		c.QueryParam("data"),
		c.QueryParam("ip"),
		c.Get("token").(string),
	)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Error inspecting domain")
	}
	return c.JSON(http.StatusOK, res)
}
