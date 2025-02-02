package HTTPServer

import (
	CustomErrors "github.com/aerosystems/checkmail-service/internal/common/custom_errors"
	"github.com/aerosystems/checkmail-service/internal/common/validators"
	"github.com/labstack/echo/v4"
	"net/http"
)

type DomainReviewRequest struct {
	Name string `json:"name" example:"gmail.com"`
	Type string `json:"type" example:"whitelist"`
}

func (r *DomainReviewRequest) Validate() error {
	if err := validators.ValidateDomainTypes(r.Type); err != nil {
		return err
	}
	if err := validators.ValidateDomainName(r.Name); err != nil {
		return err
	}
	return nil
}

// CreateReview godoc
// @Summary create top domain
// @Tags topDomains
// @Accept  json
// @Produce application/json
// @Param comment body CreateDomainRequest true "raw request body"
// @Success 201 {object} Review
// @Failure 400 {object} echo.HTTPError
// @Failure 409 {object} echo.HTTPError
// @Failure 422 {object} echo.HTTPError
// @Failure 500 {object} echo.HTTPError
// @Router /v1/reviews [post]
func (rh *ReviewHandler) CreateReview(c echo.Context) error {
	var requestPayload DomainReviewRequest
	if err := c.Bind(&requestPayload); err != nil {
		return CustomErrors.ErrInvalidRequestBody
	}
	if err := requestPayload.Validate(); err != nil {
		return err
	}
	review, err := rh.reviewUsecase.CreateReview(c.Request().Context(), requestPayload.Name, requestPayload.Type)
	if err != nil {
		return err
	}
	return c.JSON(http.StatusCreated, ModelToReview(review))
}
