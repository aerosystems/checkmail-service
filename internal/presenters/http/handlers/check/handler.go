package check

import "github.com/aerosystems/checkmail-service/internal/presenters/http/handlers"

type Handler struct {
	*handlers.BaseHandler
	inspectUsecase handlers.InspectUsecase
}

func NewHandler(
	baseHandler *handlers.BaseHandler,
	inspectUsecase handlers.InspectUsecase,
) *Handler {
	return &Handler{
		BaseHandler:    baseHandler,
		inspectUsecase: inspectUsecase,
	}
}
