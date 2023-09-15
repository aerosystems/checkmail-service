package handlers

import (
	"errors"
	"github.com/aerosystems/checkmail-service/internal/helpers"
	"github.com/aerosystems/checkmail-service/internal/models"
	"gorm.io/gorm"
	"net/http"
	"strings"
)

func (h *BaseHandler) CreateFilter(w http.ResponseWriter, r *http.Request) {
	// TODO: implement filters for customers
	_ = WriteResponse(w, http.StatusNotImplemented, NewErrorPayload(501201, "not implemented", nil))
	return
}

// CreateFilterReview godoc
// @Summary create top domain
// @Tags topDomains
// @Accept  json
// @Produce application/json
// @Param X-Api-Key header string true "api key"
// @Param comment body DomainRequest true "raw request body"
// @Success 201 {object} Response{data=models.Filter}
// @Failure 400 {object} ErrorResponse
// @Failure 409 {object} ErrorResponse
// @Failure 422 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /v1/filter [post]
func (h *BaseHandler) CreateFilterReview(w http.ResponseWriter, r *http.Request) {
	xApiKey := r.Header.Get("X-Api-Key")
	var requestPayload DomainRequest
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
		err := errors.New("root domain does not exist")
		_ = WriteResponse(w, http.StatusBadRequest, NewErrorPayload(400205, err.Error(), err))
		return
	}

	newTopDomain := models.Filter{
		Name:         requestPayload.Name,
		Type:         requestPayload.Type,
		Coverage:     requestPayload.Coverage,
		ProjectToken: xApiKey,
	}

	if err := h.topDomainRepo.Create(&newTopDomain); err != nil {
		if errors.Is(err, gorm.ErrDuplicatedKey) || strings.Contains(err.Error(), "duplicate key value violates unique constraint") {
			_ = WriteResponse(w, http.StatusConflict, NewErrorPayload(409201, "top domain already exists", err))
			return
		}
		_ = WriteResponse(w, http.StatusInternalServerError, NewErrorPayload(500202, "could not create top domain", err))
		return
	}

	_ = WriteResponse(w, http.StatusCreated, NewResponsePayload("top domain created", newTopDomain))
	return
}
