package main

import (
	"github.com/aerosystems/checkmail-service/internal/middleware"
	"github.com/aerosystems/checkmail-service/internal/presenters/rest"
)

type App struct {
	domainHandler       rest.DomainHandler
	filterHandler       rest.FilterHandler
	inspectHandler      rest.InspectHandler
	reviewHandler       rest.ReviewHandler
	oauthMiddleware     middleware.OAuthMiddleware
	basicAuthMiddleware middleware.BasicAuthMiddleware
}

func NewApp(
	domainHandler rest.DomainHandler,
	filterHandler rest.FilterHandler,
	inspectHandler rest.InspectHandler,
	reviewHandler rest.ReviewHandler,
	oauthMiddleware middleware.OAuthMiddleware,
	basicAuthMiddleware middleware.BasicAuthMiddleware,
) *App {
	return &App{
		domainHandler:       domainHandler,
		filterHandler:       filterHandler,
		inspectHandler:      inspectHandler,
		reviewHandler:       reviewHandler,
		oauthMiddleware:     oauthMiddleware,
		basicAuthMiddleware: basicAuthMiddleware,
	}
}
