package handlers

import "net/http"

func (h *BaseHandler) DomainRead(w http.ResponseWriter, r *http.Request) {
	payload := NewResponsePayload("method not implemented", nil)
	_ = WriteResponse(w, http.StatusNotImplemented, payload)
	return
}
