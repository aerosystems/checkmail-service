package RPCServer

import (
	"github.com/aerosystems/checkmail-service/internal/models"
)

type InspectUsecase interface {
	InspectData(data, clientIp, projectToken string) (*string, *models.Error)
}
