package main

import (
	"github.com/aerosystems/checkmail-service/internal/middleware"
	"github.com/aerosystems/checkmail-service/internal/presenters/rest"
	rpcServer "github.com/aerosystems/checkmail-service/internal/presenters/rpc"
	"github.com/sirupsen/logrus"
)

type App struct {
	log                 *logrus.Logger
	domainHandler       *rest.DomainHandler
	filterHandler       *rest.FilterHandler
	inspectHandler      *rest.InspectHandler
	reviewHandler       *rest.ReviewHandler
	oauthMiddleware     *middleware.OAuthMiddleware
	basicAuthMiddleware *middleware.BasicAuthMiddleware
	rpcServer           *rpcServer.RPCServer
}

func NewApp(
	log *logrus.Logger,
	domainHandler *rest.DomainHandler,
	filterHandler *rest.FilterHandler,
	inspectHandler *rest.InspectHandler,
	reviewHandler *rest.ReviewHandler,
	oauthMiddleware *middleware.OAuthMiddleware,
	basicAuthMiddleware *middleware.BasicAuthMiddleware,
	rpcServer *rpcServer.RPCServer,
) *App {
	return &App{
		log:                 log,
		domainHandler:       domainHandler,
		filterHandler:       filterHandler,
		inspectHandler:      inspectHandler,
		reviewHandler:       reviewHandler,
		oauthMiddleware:     oauthMiddleware,
		basicAuthMiddleware: basicAuthMiddleware,
		rpcServer:           rpcServer,
	}
}
