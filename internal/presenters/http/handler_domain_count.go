package HTTPServer

import (
	"github.com/aerosystems/checkmail-service/internal/models"
	"github.com/labstack/echo/v4"
	"net/http"
)

type CountResponse struct {
	Count map[string]int `json:"count"`
}

func responseFromModel(m map[models.Type]int) map[string]int {
	countMap := make(map[string]int)
	for k, v := range m {
		countMap[k.String()] = v
	}
	return countMap
}

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
	return c.JSON(http.StatusOK, CountResponse{responseFromModel(count)})
}
