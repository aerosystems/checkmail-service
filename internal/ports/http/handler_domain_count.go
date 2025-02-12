package HTTPServer

import (
	"github.com/aerosystems/checkmail-service/internal/entities"
	"github.com/labstack/echo/v4"
	"net/http"
)

type CountResponse struct {
	Whitelist int `json:"whitelist"`
	Blacklist int `json:"blacklist"`
}

func responseFromModel(m map[entities.Type]int) *CountResponse {
	var count CountResponse
	if v, ok := m[entities.BlacklistType]; ok {
		count.Blacklist = v
	}
	if v, ok := m[entities.WhitelistType]; ok {
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
