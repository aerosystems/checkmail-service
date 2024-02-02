package RPCClient

import (
	"github.com/aerosystems/checkmail-service/internal/models"
	"github.com/google/uuid"
	"net/rpc"
)

type ProjectRepo struct {
	address string // project-service:5001
}

func NewProjectRepo(address string) *ProjectRepo {
	return &ProjectRepo{
		address: address,
	}
}

func (pr *ProjectRepo) GetProjectList(userUuid uuid.UUID) (*[]models.ProjectRPCPayload, error) {
	if projectClientRPC, err := rpc.Dial("tcp", pr.address); err == nil {
		var result []models.ProjectRPCPayload
		if err := projectClientRPC.Call(
			"ProjectServer.GetProjectList",
			userUuid,
			&result,
		); err != nil {
			return nil, err
		}
		return &result, nil
	} else {
		return nil, err
	}
}

func (pr *ProjectRepo) GetProject(projectToken string) (*models.ProjectRPCPayload, error) {
	if projectClientRPC, err := rpc.Dial("tcp", pr.address); err == nil {
		var result models.ProjectRPCPayload
		if err := projectClientRPC.Call(
			"ProjectServer.GetProject",
			projectToken,
			&result,
		); err != nil {
			return nil, err
		}
		return &result, nil
	} else {
		return nil, err
	}
}
