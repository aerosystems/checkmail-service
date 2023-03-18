package handlers

import (
	"github.com/aerosystems/checkmail-service/internal/models"
	"net/http"
)

func (h *BaseHandler) DomainCreate(w http.ResponseWriter, r *http.Request) {
	var domain models.Domain

	if err := ReadRequest(w, r, &domain); err != nil {
		_ = WriteResponse(w, http.StatusBadRequest, NewErrorPayload(400201, "request payload is incorrect", err))
		return
	}

	if err := h.domainRepo.Create(&domain); err != nil {
		_ = WriteResponse(w, http.StatusInternalServerError, NewErrorPayload(500203, "could not create new Domain", err))
		return
	}

	payload := NewResponsePayload("method not implemented", nil)
	_ = WriteResponse(w, http.StatusNotImplemented, payload)
	return
}
