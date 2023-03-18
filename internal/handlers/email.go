package handlers

import "net/http"

func (h *BaseHandler) Email(w http.ResponseWriter, r *http.Request) {
	payload := NewResponsePayload("method not implemented", nil)
	_ = WriteResponse(w, http.StatusNotImplemented, payload)
	return
}
