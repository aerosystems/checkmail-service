package RPCServer

import CustomError "github.com/aerosystems/checkmail-service/pkg/custom_error"

type InspectService interface {
	InspectData(data, clientIp, projectToken string) (*string, *CustomError.Error)
}
