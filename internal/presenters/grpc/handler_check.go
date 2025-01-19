package GRPCServer

import (
	"context"
	"github.com/aerosystems/checkmail-service/internal/common/protobuf/checkmail"
)

type CheckHandler struct {
	inspectUsecase InspectUsecase
	checkmail.UnimplementedCheckmailServiceServer
}

func NewCheckHandler(inspectUsecase InspectUsecase) *CheckHandler {
	return &CheckHandler{
		inspectUsecase: inspectUsecase,
	}
}

func (h CheckHandler) Inspect(_ context.Context, req *checkmail.InspectRequest) (*checkmail.InspectResponse, error) {
	domainType, err := h.inspectUsecase.InspectData(req.Data, req.ClientIp, req.ProjectToken)
	if err != nil {
		return nil, err
	}
	return &checkmail.InspectResponse{
		DomainType: domainType.String(),
	}, nil
}
