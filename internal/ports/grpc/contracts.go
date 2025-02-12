package GRPCServer

import (
	"context"
	"github.com/aerosystems/checkmail-service/internal/entities"
)

type InspectUsecase interface {
	InspectData(ctx context.Context, data, clientIp, projectToken string) (*entities.Type, error)
}
