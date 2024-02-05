package rpcServer

import (
	"fmt"
	"github.com/labstack/gommon/log"
	"github.com/sirupsen/logrus"
	"net"
	"net/rpc"
)

const rpcPort = 5001

type RPCServer struct {
	rpcPort        int
	log            *logrus.Logger
	inspectUsecase InspectUsecase
}

func NewRPCServer(
	log *logrus.Logger,
	inspectUsecase InspectUsecase,
) *RPCServer {
	return &RPCServer{
		rpcPort:        rpcPort,
		log:            log,
		inspectUsecase: inspectUsecase,
	}
}

func (cs RPCServer) Listen() error {
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
