package RPCServer

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"net"
	"net/rpc"
)

type CheckmailServer struct {
	rpcPort        int
	log            *logrus.Logger
	InspectService *usecase.InspectService
}

func NewCheckmailServer(
	rpcPort int,
	inspectService *usecase.InspectService,
) *CheckmailServer {
	return &CheckmailServer{
		rpcPort:        rpcPort,
		InspectService: inspectService,
	}
}

func (cs *CheckmailServer) Listen() error {
	listen, err := net.Listen("tcp", fmt.Sprintf("0.0.0.0:%d", cs.rpcPort))
	if err != nil {
		return err
	}
	defer listen.Close()

	for {
		rpcConn, err := listen.Accept()
		if err != nil {
			continue
		}
		go rpc.ServeConn(rpcConn)
	}
}
