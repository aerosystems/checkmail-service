package check

import (
	"fmt"
	"github.com/labstack/echo/v4"
	"net/http"
	"time"
)

type InspectRequest struct {
	Data     string `json:"data"`
	ClientIp string `json:"clientIp,omitempty"`
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
// @Router /v1/inspect [post]
func (ch Handler) Inspect(c echo.Context) error {
	start := time.Now()
	var requestPayload InspectRequest
	if err := c.Bind(&requestPayload); err != nil {
		return ch.ErrorResponse(c, http.StatusUnprocessableEntity, "could not read request body", err)
	}
	domainType, err := ch.inspectUsecase.InspectData(requestPayload.Data, requestPayload.ClientIp, c.Request().Header.Get("X-Api-Key"))
	if err != nil {
		return ch.ErrorResponse(c, http.StatusBadRequest, err.Message, err.Error())
	}
	duration := time.Since(start)
	return ch.SuccessResponse(c, http.StatusOK, fmt.Sprintf("%s is defined as %s per %d milliseconds", requestPayload.Data, *domainType, duration.Milliseconds()), *domainType)
}
