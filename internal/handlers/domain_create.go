package handlers

import (
	"errors"
	"fmt"
	"github.com/aerosystems/checkmail-service/internal/models"
	"github.com/aerosystems/checkmail-service/pkg/validators"
	"gorm.io/gorm"
	"net/http"
	"strings"
)

type CreateDomainRequest struct {
	Name     string `json:"name" example:"gmail.com"`
	Type     string `json:"type" example:"whitelist"`
	Coverage string `json:"coverage" example:"equals"`
}

// CreateDomain godoc
// @Summary create domain
// @Tags domains
// @Accept  json
// @Produce application/json
// @Param comment body CreateDomainRequest true "raw request body"
// @Security BearerAuth
// @Success 200 {object} Response{data=models.Domain}
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Failure 403 {object} ErrorResponse
// @Failure 409 {object} ErrorResponse
// @Failure 422 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /v1/domains [post]
func (h *BaseHandler) CreateDomain(w http.ResponseWriter, r *http.Request) {
	var requestPayload CreateDomainRequest

	if err := ReadRequest(w, r, &requestPayload); err != nil {
		_ = WriteResponse(w, http.StatusUnprocessableEntity, NewErrorPayload(422201, "could not read request body", err))
		return
	}

	if requestPayload.Name == "" {
		err := errors.New("name does not exists or empty")
		_ = WriteResponse(w, http.StatusUnprocessableEntity, NewErrorPayload(422204, err.Error(), err))
		return
	}

	if requestPayload.Type == "" {
		err := errors.New("type does not exists or empty")
		_ = WriteResponse(w, http.StatusUnprocessableEntity, NewErrorPayload(422205, err.Error(), err))
		return
	}

	if requestPayload.Coverage == "" {
		err := errors.New("coverage does not exists or empty")
		_ = WriteResponse(w, http.StatusUnprocessableEntity, NewErrorPayload(422206, err.Error(), err))
		return
	}

	if err := validators.ValidateDomainTypes(requestPayload.Type); err != nil {
		_ = WriteResponse(w, http.StatusBadRequest, NewErrorPayload(400203, err.Error(), err))
		return
	}

	if err := validators.ValidateDomainCoverages(requestPayload.Coverage); err != nil {
		_ = WriteResponse(w, http.StatusBadRequest, NewErrorPayload(400203, err.Error(), err))
		return
	}

	if err := validators.ValidateDomain(requestPayload.Name); err != nil {
		_ = WriteResponse(w, http.StatusBadRequest, NewErrorPayload(400204, err.Error(), err))
		return
	}

	arrDomain := strings.Split(requestPayload.Name, ".")
	root := arrDomain[len(arrDomain)-1]
	rootDomain, _ := h.rootDomainRepo.FindByName(root)
	if rootDomain == nil {
		err := errors.New("root domain does not exist")
		_ = WriteResponse(w, http.StatusBadRequest, NewErrorPayload(400205, err.Error(), err))
		return
	}

	newDomain := models.Domain{
		Name:     requestPayload.Name,
		Type:     requestPayload.Type,
		Coverage: requestPayload.Coverage,
	}

	if err := h.domainRepo.Create(&newDomain); err != nil {
		if errors.Is(err, gorm.ErrDuplicatedKey) || strings.Contains(err.Error(), "duplicate key value violates unique constraint") {
			_ = WriteResponse(w, http.StatusConflict, NewErrorPayload(409202, fmt.Sprintf("domain %s already exists", newDomain.Name), err))
			return
		}
		_ = WriteResponse(w, http.StatusInternalServerError, NewErrorPayload(500203, "could not create new Domain", err))
		return
	}

	_ = WriteResponse(w, http.StatusOK, NewResponsePayload("domain successfully created", newDomain))
	return
}
