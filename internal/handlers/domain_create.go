package handlers

import (
	"errors"
	"fmt"
	"github.com/aerosystems/checkmail-service/internal/helpers"
	"github.com/aerosystems/checkmail-service/internal/models"
	"gorm.io/gorm"
	"net/http"
)

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

	if err := helpers.ValidateDomainTypes(newDomain.Type); err != nil {
		_ = WriteResponse(w, http.StatusBadRequest, NewErrorPayload(400203, err.Error(), err))
		return
	}

	if err := helpers.ValidateDomainCoverages(newDomain.Coverage); err != nil {
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
