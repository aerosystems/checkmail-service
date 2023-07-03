package handlers

import (
	"errors"
	"fmt"
	"github.com/aerosystems/checkmail-service/internal/models"
	"github.com/aerosystems/checkmail-service/pkg/validators"
	"gorm.io/gorm"
	"net/http"
)

// DomainCreate godoc
// @Summary create domain
// @Tags domains
// @Accept  json
// @Produce application/json
// @Param comment body models.Domain true "raw request body"
// @Param Authorization header string true "should contain Access Token, with the Bearer started"
// @Success 200 {object} Response{data=models.Domain}
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Failure 403 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /domains [post]
func (h *BaseHandler) DomainCreate(w http.ResponseWriter, r *http.Request) {
	var newDomain models.Domain

	if err := ReadRequest(w, r, &newDomain); err != nil {
		_ = WriteResponse(w, http.StatusBadRequest, NewErrorPayload(400201, "request payload is incorrect", err))
		return
	}

	if newDomain.Name == "" {
		err := errors.New("claim Name does not exists or empty")
		_ = WriteResponse(w, http.StatusBadRequest, NewErrorPayload(400204, err.Error(), err))
		return
	}

	if newDomain.Type == "" {
		err := errors.New("claim Name does not exists or empty")
		_ = WriteResponse(w, http.StatusBadRequest, NewErrorPayload(400205, err.Error(), err))
		return
	}

	if newDomain.Coverage == "" {
		err := errors.New("claim Name does not exists or empty")
		_ = WriteResponse(w, http.StatusBadRequest, NewErrorPayload(400206, err.Error(), err))
		return
	}

	if err := validators.ValidateDomainTypes(newDomain.Type); err != nil {
		_ = WriteResponse(w, http.StatusBadRequest, NewErrorPayload(400203, err.Error(), err))
		return
	}

	if err := validators.ValidateDomainCoverages(newDomain.Coverage); err != nil {
		_ = WriteResponse(w, http.StatusBadRequest, NewErrorPayload(400203, err.Error(), err))
		return
	}

	if err := h.domainRepo.Create(&newDomain); err != nil {
		if err == gorm.ErrDuplicatedKey {
			_ = WriteResponse(w, http.StatusInternalServerError, NewErrorPayload(500202, fmt.Sprintf("domain %s already exists", newDomain.Name), err))
			return
		}
		_ = WriteResponse(w, http.StatusInternalServerError, NewErrorPayload(500203, "could not create new Domain", err))
		return
	}

	payload := NewResponsePayload("domain successfully created", newDomain)
	_ = WriteResponse(w, http.StatusOK, payload)
	return
}
