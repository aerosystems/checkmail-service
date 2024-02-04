package rest

import (
	"github.com/aerosystems/checkmail-service/internal/validators"
	CustomError "github.com/aerosystems/checkmail-service/pkg/custom_error"
	"github.com/labstack/echo/v4"
	"net/http"
)

type ReviewHandler struct {
	*BaseHandler
	reviewUsecase ReviewUsecase
}

func NewReviewHandler(baseHandler *BaseHandler, reviewUsecase ReviewUsecase) *ReviewHandler {
	return &ReviewHandler{
		BaseHandler:   baseHandler,
		reviewUsecase: reviewUsecase,
	}
}

type DomainReviewRequest struct {
	Name string `json:"name" example:"gmail.com"`
	Type string `json:"type" example:"whitelist"`
}

func (r *DomainReviewRequest) Validate() *CustomError.Error {
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
// @Success 201 {object} Response{data=models.DomainReview}
// @Failure 400 {object} ErrorResponse
// @Failure 409 {object} ErrorResponse
// @Failure 422 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /v1/reviews [post]
func (rh *ReviewHandler) CreateReview(c echo.Context) error {
	var requestPayload DomainReviewRequest
	if err := c.Bind(&requestPayload); err != nil {
		return rh.ErrorResponse(c, http.StatusUnprocessableEntity, "could not read request body", err)
	}
	if err := requestPayload.Validate(); err != nil {
		return rh.ErrorResponse(c, http.StatusBadRequest, err.Message, err.Error())
	}
	review, err := rh.reviewUsecase.CreateReview(requestPayload.Name, requestPayload.Type)
	if err != nil {
		return rh.ErrorResponse(c, http.StatusInternalServerError, "could not create review", err)
	}
	return rh.SuccessResponse(c, http.StatusCreated, "review successfully created", review)
}
