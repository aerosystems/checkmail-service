package handlers

import "net/http"

func (h *BaseHandler) DomainCreate(w http.ResponseWriter, r *http.Request) {
	payload := NewResponsePayload("method not implemented", nil)
	_ = WriteResponse(w, http.StatusNotImplemented, payload)
	return
}
