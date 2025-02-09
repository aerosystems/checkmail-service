package HTTPServer

import (
	"github.com/aerosystems/checkmail-service/internal/models"
	"github.com/labstack/echo/v4"
	"net/http"
)

type CountResponse struct {
	Whitelist int `json:"whitelist"`
	Blacklist int `json:"blacklist"`
}

func responseFromModel(m map[models.Type]int) *CountResponse {
	var count CountResponse
	if v, ok := m[models.BlacklistType]; ok {
		count.Blacklist = v
	}
	if v, ok := m[models.WhitelistType]; ok {
		count.Whitelist = v
	}
	return &count
}

// Count godoc
// @Summary count Domains
// @Tags domains
// @Accept  json
// @Produce application/json
// @Success 200 {object} CountResponse
// @Failure 500 {object} echo.HTTPError
// @Router /v1/domains/count [post]
func (h Handler) Count(c echo.Context) error {
	count, err := h.domainUsecase.CountDomains(c.Request().Context())
	if err != nil {
		return err
	}
	return c.JSON(http.StatusOK, responseFromModel(count))
}
