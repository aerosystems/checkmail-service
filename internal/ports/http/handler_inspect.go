package HTTPServer

import (
	"fmt"
	"github.com/aerosystems/checkmail-service/internal/entities"
	"github.com/labstack/echo/v4"
	"net/http"
	"time"
)

type InspectRequest struct {
	Data     string `json:"data"`
	ClientIp string `json:"clientIp,omitempty"`
}

type InspectRequestPublic struct {
	Data string `json:"data"`
}

type InspectResponse struct {
	Message string `json:"message"`
	Domain  Domain `json:"domain"`
}

// Inspect godoc
// @Summary get information about domain name or email address
// @Tags inspect
// @Accept  json
// @Produce application/json
// @Param X-Api-Key header string true "api key"
// @Param data body InspectRequest true "raw request body"
// @Success 200 {object} Response{data=string}
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /v1/data/inspect [post]
func (h Handler) Inspect(c echo.Context) error {
	start := time.Now()
	var requestPayload InspectRequest
	if err := c.Bind(&requestPayload); err != nil {
		return entities.ErrInvalidRequestBody
	}
	domainName, domainType, err := h.inspectUsecase.InspectDataWithAuth(c.Request().Context(), requestPayload.Data,
		requestPayload.ClientIp, getAPIKeyFromContext(c))
	if err != nil {
		return err
	}
	duration := time.Since(start)
	return c.JSON(http.StatusOK, InspectResponse{Message: fmt.Sprintf("%s is defined as %s per %d ms", requestPayload.Data, domainType.String(), duration.Milliseconds()),
		Domain: Domain{
			Name: domainName,
			Type: domainType.String(),
		},
	})
}

func (h Handler) InspectPublic(c echo.Context) error {
	start := time.Now()
	var requestPayload InspectRequestPublic
	if err := c.Bind(&requestPayload); err != nil {
		return entities.ErrInvalidRequestBody
	}
	domainName, domainType, err := h.inspectUsecase.InspectData(c.Request().Context(), requestPayload.Data)
	if err != nil {
		return err
	}
	duration := time.Since(start)
	return c.JSON(http.StatusOK, InspectResponse{Message: fmt.Sprintf("%s is defined as %s per %d milliseconds", requestPayload.Data, domainType.String(), duration.Milliseconds()),
		Domain: Domain{
			Name: domainName,
			Type: domainType.String(),
		},
	})
}
