package main

import (
	"github.com/aerosystems/checkmail-service/internal/config"
	"github.com/aerosystems/checkmail-service/internal/http"
	"github.com/aerosystems/checkmail-service/internal/presenters/rpc"
	"github.com/sirupsen/logrus"
)

type App struct {
	log        *logrus.Logger
	cfg        *config.Config
	httpServer *HTTPServer.Server
	rpcServer  *RPCServer.Server
}

func NewApp(
	log *logrus.Logger,
	cfg *config.Config,
	httpServer *HTTPServer.Server,
	rpcServer *RPCServer.Server,
) *App {
	return &App{
		log:        log,
		cfg:        cfg,
		httpServer: httpServer,
		rpcServer:  rpcServer,
	}
}
