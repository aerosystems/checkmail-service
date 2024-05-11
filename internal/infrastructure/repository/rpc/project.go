package RpcRepo

import (
	"github.com/aerosystems/checkmail-service/internal/models"
	"github.com/aerosystems/checkmail-service/pkg/rpc_client"
	"github.com/google/uuid"
)

type ProjectRepo struct {
	rpcClient *RpcClient.ReconnectRpcClient
}

func NewProjectRepo(rpcClient *RpcClient.ReconnectRpcClient) *ProjectRepo {
	return &ProjectRepo{
		rpcClient: rpcClient,
	}
}

func (pr *ProjectRepo) GetProjectList(userUuid uuid.UUID) (*[]models.ProjectRPCPayload, error) {
	var result []models.ProjectRPCPayload
	if err := pr.rpcClient.Call(
		"Server.GetProjectList",
		userUuid,
		&result,
	); err != nil {
		return nil, err
	}
	return &result, nil
}

func (pr *ProjectRepo) GetProject(projectToken string) (*models.ProjectRPCPayload, error) {
	var result models.ProjectRPCPayload
	if err := pr.rpcClient.Call(
		"Server.GetProject",
		projectToken,
		&result,
	); err != nil {
		return nil, err
	}
	return &result, nil
}
