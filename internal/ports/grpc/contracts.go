package GRPCServer

import (
	"context"
	"github.com/aerosystems/checkmail-service/internal/entities"
)

type InspectUsecase interface {
	InspectDataWithAuth(ctx context.Context, data, clientIp, projectToken string) (string, *entities.Type, error)
}
