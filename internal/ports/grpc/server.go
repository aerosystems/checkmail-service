package GRPCServer

import (
	"github.com/aerosystems/common-service/gen/protobuf/checkmail"
	"github.com/aerosystems/common-service/presenters/grpcserver"
	"github.com/sirupsen/logrus"
)

type Server struct {
	grpcServer *grpcserver.Server
}

func NewGRPCServer(
	cfg *grpcserver.Config,
	log *logrus.Logger,
	checkHandler *CheckHandler,
) *Server {
	server := grpcserver.NewGRPCServer(cfg, log)

	server.RegisterService(checkmail.CheckmailService_ServiceDesc, checkHandler)

	return &Server{
		grpcServer: server,
	}
}

func (s *Server) Run() error {
	return s.grpcServer.Run()
}
