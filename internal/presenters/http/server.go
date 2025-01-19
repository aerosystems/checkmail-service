package HTTPServer

import (
	"fmt"
	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"
)

type Server struct {
	echo                   *echo.Echo
	port                   int
	log                    *logrus.Logger
	firebaseAuthMiddleware *FirebaseAuth
	apiKeyAuthMiddleware   *ApiKeyAuth
	domainHandler          *DomainHandler
	filterHandler          *FilterHandler
	checkHandler           *CheckHandler
	reviewHandler          *ReviewHandler
	accessHandler          *AccessHandler
}

type Handlers struct {
	DomainHandler *DomainHandler
	FilterHandler *FilterHandler
	CheckHandler  *CheckHandler
	ReviewHandler *ReviewHandler
	AccessHandler *AccessHandler
}

type Middlewares struct {
	FirebaseAuthMiddleware *FirebaseAuth
	ApiKeyAuthMiddleware   *ApiKeyAuth
}

func NewServer(
	port int,
	log *logrus.Logger,
	errorHandler *echo.HTTPErrorHandler,
	handlers Handlers,
	middlewares Middlewares,
) *Server {
	server := &Server{
		port:                   port,
		log:                    log,
		echo:                   echo.New(),
		firebaseAuthMiddleware: middlewares.FirebaseAuthMiddleware,
		apiKeyAuthMiddleware:   middlewares.ApiKeyAuthMiddleware,
		domainHandler:          handlers.DomainHandler,
		filterHandler:          handlers.FilterHandler,
		checkHandler:           handlers.CheckHandler,
		reviewHandler:          handlers.ReviewHandler,
		accessHandler:          handlers.AccessHandler,
	}
	if errorHandler != nil {
		server.echo.HTTPErrorHandler = *errorHandler
	}
	return server
}

func (s *Server) Run() error {
	s.setupMiddleware()
	s.setupRoutes()
	return s.echo.Start(fmt.Sprintf(":%d", s.port))
}
