package main

import (
	"github.com/aerosystems/checkmail-service/internal/common/config"
	GRPCServer "github.com/aerosystems/checkmail-service/internal/ports/grpc"
	HTTPServer "github.com/aerosystems/checkmail-service/internal/ports/http"
	"github.com/sirupsen/logrus"
)

type App struct {
	log        *logrus.Logger
	cfg        *config.Config
	httpServer *HTTPServer.Server
	grpcServer *GRPCServer.Server
}

func NewApp(
	log *logrus.Logger,
	cfg *config.Config,
	httpServer *HTTPServer.Server,
	grpcServer *GRPCServer.Server,
) *App {
	return &App{
		log:        log,
		cfg:        cfg,
		httpServer: httpServer,
		grpcServer: grpcServer,
	}
}
