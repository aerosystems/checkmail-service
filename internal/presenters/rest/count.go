package rest

import (
	"net/http"
)

// Count godoc
// @Summary count Domains
// @Tags domains
// @Accept  json
// @Produce application/json
// @Success 200 {object} Response
// @Failure 500 {object} ErrorResponse
// @Router /v1/domains/count [post]
func (h *BaseHandler) Count(w http.ResponseWriter, r *http.Request) {
	count, err := h.domainRepo.Count()
	if err != nil {
		_ = WriteResponse(w, http.StatusInternalServerError, NewErrorPayload(500201, "could not count Domains", err))
		return
	}
	_ = WriteResponse(w, http.StatusOK, NewResponsePayload("domain successfully counted", count))
	return
}
