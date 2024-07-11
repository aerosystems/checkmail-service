package access

import "github.com/aerosystems/checkmail-service/internal/presenters/http/handlers"

type Handler struct {
	accessUsecase handlers.AccessUsecase
}

func NewHandler(accessUsecase handlers.AccessUsecase) *Handler {
	return &Handler{accessUsecase: accessUsecase}
}
