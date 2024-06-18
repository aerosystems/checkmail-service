package domain

import (
	"errors"
	CustomErrors "github.com/aerosystems/checkmail-service/internal/common/custom_errors"
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
func (dh Handler) Count(c echo.Context) error {
	count, err := dh.domainUsecase.CountDomains()
	if err != nil {
		var apiErr CustomErrors.ApiError
		if errors.As(err, &apiErr) {
			return dh.ErrorResponse(c, apiErr.HttpCode, apiErr.Message, err)
		}
		return dh.ErrorResponse(c, CustomErrors.ErrDomainInternalCount.HttpCode, CustomErrors.ErrDomainInternalCount.Message, err)
	}
	return dh.SuccessResponse(c, http.StatusOK, "domains successfully counted", count)
}
