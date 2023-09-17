package RPCClient

import (
	"net/rpc"
	"time"
)

type ProjectRPCPayload struct {
	ID         int
	UserID     int
	Name       string
	Token      string
	AccessTime time.Time
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
