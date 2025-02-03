package GRPCServer

import (
	"context"
	"github.com/aerosystems/common-service/gen/protobuf/checkmail"
)

type CheckService struct {
	inspectUsecase InspectUsecase
	checkmail.UnimplementedCheckmailServiceServer
}

func NewCheckService(inspectUsecase InspectUsecase) *CheckService {
	return &CheckService{
		inspectUsecase: inspectUsecase,
	}
}

func (cs CheckService) Inspect(ctx context.Context, req *checkmail.InspectRequest) (*checkmail.InspectResponse, error) {
	domainType, err := cs.inspectUsecase.InspectData(ctx, req.Data, req.ClientIp, req.ProjectToken)
	if err != nil {
		return nil, err
	}
	return &checkmail.InspectResponse{
		DomainType: domainType.String(),
	}, nil
}
