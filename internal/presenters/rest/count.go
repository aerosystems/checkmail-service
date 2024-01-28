package rest

import (
	"github.com/labstack/echo/v4"
	"net/http"
)

// Count godoc
// @Summary count Domains
// @Tags domains
// @Accept  json
// @Produce application/json
// @Success 200 {object} Response
// @Failure 500 {object} ErrorResponse
// @Router /v1/domains/count [post]
func (h *BaseHandler) Count(c echo.Context) error {
	count, err := h.domainRepo.Count()
	if err != nil {
		return h.ErrorResponse(c, http.StatusInternalServerError, "could not count Domains", err)

	}
	return h.SuccessResponse(c, http.StatusOK, "domains successfully counted", count)
}
