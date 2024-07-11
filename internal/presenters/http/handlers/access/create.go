package access

import (
	"encoding/json"
	CustomErrors "github.com/aerosystems/checkmail-service/internal/common/custom_errors"
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
	AccessTime       time.Time `json:"accessTime"`
}

func ModelToCreateAccessEvent(access *models.Access) CreateAccessEvent {
	return CreateAccessEvent{
		Token:            access.Token,
		SubscriptionType: access.SubscriptionType.String(),
		AccessTime:       access.AccessTime,
	}
}

func (h Handler) CreateAccess(c echo.Context) error {
	var req CreateAccessRequest
	if err := c.Bind(&req); err != nil {
		return CustomErrors.ErrInvalidRequestBody
	}
	var event CreateAccessEvent
	if err := json.Unmarshal(req.Message.Data, &event); err != nil {
		return CustomErrors.ErrInvalidRequestPayload
	}
	access, err := h.accessUsecase.CreateAccess(event.Token, event.SubscriptionType, event.AccessTime)
	if err != nil {
		return err
	}
	return c.JSON(http.StatusOK, ModelToCreateAccessEvent(access))
}
