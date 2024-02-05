package RPCServer

import (
	"fmt"
	"github.com/labstack/gommon/log"
	"github.com/sirupsen/logrus"
	"net"
	"net/rpc"
)

const rpcPort = 5001

type CheckmailServer struct {
	rpcPort        int
	log            *logrus.Logger
	InspectService InspectService
}

func NewCheckmailServer(
	inspectService InspectService,
) *CheckmailServer {
	return &CheckmailServer{
		rpcPort:        rpcPort,
		InspectService: inspectService,
	}
}

func (cs *CheckmailServer) Listen() error {
	log.Infof("starting checkmail-service RPC server on port %d\n", rpcPort)
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
