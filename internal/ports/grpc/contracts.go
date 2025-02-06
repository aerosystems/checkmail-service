package GRPCServer

import (
	"context"
	"github.com/aerosystems/checkmail-service/internal/models"
)

type InspectUsecase interface {
	InspectData(ctx context.Context, data, clientIp, projectToken string) (models.Type, error)
}
