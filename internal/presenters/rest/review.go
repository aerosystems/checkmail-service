package rest

import (
	"errors"
	"github.com/aerosystems/checkmail-service/internal/helpers"
	"github.com/aerosystems/checkmail-service/internal/models"
	CustomError "github.com/aerosystems/checkmail-service/pkg/custom_error"
	"github.com/aerosystems/checkmail-service/pkg/validators"
	"github.com/labstack/echo/v4"
	"net/http"
)

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

// CreateDomainReview godoc
// @Summary create top domain
// @Tags topDomains
// @Accept  json
// @Produce application/json
// @Param comment body DomainRequest true "raw request body"
// @Success 201 {object} Response{data=models.Filter}
// @Failure 400 {object} ErrorResponse
// @Failure 409 {object} ErrorResponse
// @Failure 422 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /v1/filter [post]
func (h *BaseHandler) CreateDomainReview(c echo.Context) error {
	var requestPayload DomainReviewRequest
	if err := c.Bind(&requestPayload); err != nil {
		return h.ErrorResponse(c, http.StatusUnprocessableEntity, "could not read request body", err)
	}
	if err := requestPayload.Validate(); err != nil {
		return h.ErrorResponse(c, http.StatusBadRequest, err.Message, err.Error())
	}

	root, _ := helpers.GetRootDomain(requestPayload.Name)
	rootDomain, _ := h.rootDomainRepo.FindByName(root)
	if rootDomain == nil {
		err := errors.New("domain does not exist")
		return h.ErrorResponse(c, http.StatusBadRequest, err.Error(), err)
	}

	newDomainReview := models.DomainReview{
		Name: requestPayload.Name,
		Type: requestPayload.Type,
	}
	if err := h.domainReviewRepo.Create(&newDomainReview); err != nil {
		return h.ErrorResponse(c, http.StatusInternalServerError, "could not create domain review", err)
	}
	return h.SuccessResponse(c, http.StatusCreated, "domain review successfully created", newDomainReview)
}
