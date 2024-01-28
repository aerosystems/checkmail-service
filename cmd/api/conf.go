package main

import (
	"github.com/aerosystems/checkmail-service/internal/middleware"
	"github.com/aerosystems/checkmail-service/internal/presenters/rest"
)

type Config struct {
	baseHandler         *rest.BaseHandler
	oauthMiddleware     middleware.OAuthMiddleware
	basicAuthMiddleware middleware.BasicAuthMiddleware
}

func NewConfig(baseHandler *rest.BaseHandler, oauthMiddleware middleware.OAuthMiddleware, basicAuthMiddleware middleware.BasicAuthMiddleware) *Config {
	return &Config{
		baseHandler:         baseHandler,
		oauthMiddleware:     oauthMiddleware,
		basicAuthMiddleware: basicAuthMiddleware,
	}
}
