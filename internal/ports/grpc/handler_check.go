package GRPCServer

import (
	"context"
	"github.com/aerosystems/common-service/gen/protobuf/checkmail"
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

func (h CheckHandler) Inspect(ctx context.Context, req *checkmail.InspectRequest) (*checkmail.InspectResponse, error) {
	domainType, err := h.inspectUsecase.InspectData(ctx, req.Data, req.ClientIp, req.ProjectToken)
	if err != nil {
		return nil, err
	}
	return &checkmail.InspectResponse{
		DomainType: domainType.String(),
	}, nil
}
