package HTTPServer

import (
	"encoding/json"
	"github.com/aerosystems/checkmail-service/internal/models"
	"github.com/labstack/echo/v4"
	"net/http"
	"time"
)

type CreateAccessRequest struct {
	CreateAccessRequestBody
}

type CreateAccessRequestBody struct {
	Message struct {
		Data []byte `json:"data"`
	} `json:"message"`
	Subscription string `json:"subscription"`
}

type CreateAccessEvent struct {
	Token            string    `json:"token"`
	SubscriptionType string    `json:"subscriptionType"`
	AccessCount      int       `json:"accessCount"`
	AccessTime       time.Time `json:"accessTime"`
}

func (h Handler) CreateAccess(c echo.Context) error {
	var req CreateAccessRequest
	if err := c.Bind(&req); err != nil {
		return models.ErrInvalidRequestBody
	}
	var event CreateAccessEvent
	if err := json.Unmarshal(req.Message.Data, &event); err != nil {
		return models.ErrInvalidRequestPayload
	}
	if err := h.accessUsecase.CreateAccess(c.Request().Context(), event.Token, event.SubscriptionType, event.AccessCount, event.AccessTime); err != nil {
		return err
	}
	return c.NoContent(http.StatusCreated)
}
