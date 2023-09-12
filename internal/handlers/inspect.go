package handlers

import (
	"fmt"
	"net/http"
	"time"
)

type InspectRequestPayload struct {
	Data     string `json:"data"`
	ClientIp string `json:"clientIp,omitempty"`
}

// Inspect godoc
// @Summary get information about domain name or email address
// @Tags inspect
// @Accept  json
// @Produce application/json
// @Param data body InspectRequestPayload true "raw request body"
// @Security X-Api-Key
// @Success 200 {object} Response{data=string}
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /v1/inspect [post]
func (h *BaseHandler) Inspect(w http.ResponseWriter, r *http.Request) {
	start := time.Now()
	var requestPayload InspectRequestPayload
	if err := ReadRequest(w, r, &requestPayload); err != nil {
		_ = WriteResponse(w, http.StatusUnprocessableEntity, NewErrorPayload(422001, "could not read request body", err))
		return
	}
	duration := time.Since(start)
	domainType, err := h.inspectService.InspectData(requestPayload.Data, requestPayload.ClientIp)
	if err != nil {
		_ = WriteResponse(w, http.StatusBadRequest, NewErrorPayload(err.Code, err.Message, err.Error()))
		return
	}
	_ = WriteResponse(w, http.StatusOK, NewResponsePayload(fmt.Sprintf("%s is defined as %s per %d milliseconds", requestPayload.Data, *domainType, duration.Milliseconds()), *domainType))
	return
}
