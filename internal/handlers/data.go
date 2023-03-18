package handlers

import "net/http"

func (h *BaseHandler) Data(w http.ResponseWriter, r *http.Request) {
	payload := NewResponsePayload("method not implemented", nil)
	_ = WriteResponse(w, http.StatusNotImplemented, payload)
	return
}
