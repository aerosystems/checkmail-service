package HttpServer

import (
	"github.com/labstack/echo/v4"
	"net/http"
)

type CountResponse struct {
	CountMap
}

type CountMap map[string]int

// Count godoc
// @Summary count Domains
// @Tags domains
// @Accept  json
// @Produce application/json
// @Success 200 {object} CountResponse
// @Failure 500 {object} echo.HTTPError
// @Router /v1/domains/count [post]
func (dh DomainHandler) Count(c echo.Context) error {
	count, err := dh.domainUsecase.CountDomains()
	if err != nil {
		return err
	}
	return c.JSON(http.StatusOK, CountResponse{count})
}
