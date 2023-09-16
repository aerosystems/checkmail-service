package handlers

import (
	"errors"
	"github.com/aerosystems/checkmail-service/internal/helpers"
	"github.com/aerosystems/checkmail-service/internal/models"
	CustomError "github.com/aerosystems/checkmail-service/pkg/custom_error"
	"github.com/aerosystems/checkmail-service/pkg/validators"
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
func (h *BaseHandler) CreateDomainReview(w http.ResponseWriter, r *http.Request) {
	var requestPayload DomainReviewRequest
	if err := ReadRequest(w, r, &requestPayload); err != nil {
		WriteResponse(w, http.StatusUnprocessableEntity, NewErrorPayload(422201, "could not read request body", err))
		return
	}

	if err := requestPayload.Validate(); err != nil {
		WriteResponse(w, http.StatusBadRequest, NewErrorPayload(err.Code, err.Message, err.Error()))
		return
	}

	root, _ := helpers.GetRootDomain(requestPayload.Name)
	rootDomain, _ := h.rootDomainRepo.FindByName(root)
	if rootDomain == nil {
		err := errors.New("domain does not exist")
		_ = WriteResponse(w, http.StatusBadRequest, NewErrorPayload(400205, err.Error(), err))
		return
	}

	newDomainReview := models.DomainReview{
		Name: requestPayload.Name,
		Type: requestPayload.Type,
	}

	if err := h.domainReviewRepo.Create(&newDomainReview); err != nil {
		_ = WriteResponse(w, http.StatusInternalServerError, NewErrorPayload(500202, "could not create domain review", err))
		return
	}

	_ = WriteResponse(w, http.StatusCreated, NewResponsePayload("domain review created", newDomainReview))
	return
}
