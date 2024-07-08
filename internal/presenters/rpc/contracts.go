package RpcServer

import (
	"github.com/aerosystems/checkmail-service/internal/models"
)

type InspectUsecase interface {
	InspectData(string, string, string) (models.DomainType, error)
}
