package RPCClient

import (
	"github.com/google/uuid"
	"net/rpc"
)

type ProjectRPCPayload struct {
	Id       int
	UserUuid uuid.UUID
	Name     string
	Token    string
}

func GetProjectList(userId int) (*[]ProjectRPCPayload, error) {
	if projectClientRPC, err := rpc.Dial("tcp", "project-service:5001"); err == nil {
		var result []ProjectRPCPayload
		if err := projectClientRPC.Call(
			"ProjectServer.GetProjectList",
			userId,
			&result,
		); err != nil {
			return nil, err
		}
		return &result, nil
	} else {
		return nil, err
	}
}

func GetProject(projectToken string) (*ProjectRPCPayload, error) {
	if projectClientRPC, err := rpc.Dial("tcp", "project-service:5001"); err == nil {
		var result ProjectRPCPayload
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
