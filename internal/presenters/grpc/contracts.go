package GRPCServer

import (
	"github.com/aerosystems/checkmail-service/internal/models"
)

type InspectUsecase interface {
	InspectData(data, clientIp, projectToken string) (models.Type, error)
}
