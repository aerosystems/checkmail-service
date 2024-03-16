package main

import (
	"github.com/aerosystems/checkmail-service/internal/config"
	"github.com/aerosystems/checkmail-service/internal/infrastructure/http"
	"github.com/aerosystems/checkmail-service/internal/infrastructure/rpc"
	"github.com/sirupsen/logrus"
)

type App struct {
	log        *logrus.Logger
	cfg        *config.Config
	httpServer *HttpServer.Server
	rpcServer  *RpcServer.Server
}

func NewApp(
	log *logrus.Logger,
	cfg *config.Config,
	httpServer *HttpServer.Server,
	rpcServer *RpcServer.Server,
) *App {
	return &App{
		log:        log,
		cfg:        cfg,
		httpServer: httpServer,
		rpcServer:  rpcServer,
	}
}
