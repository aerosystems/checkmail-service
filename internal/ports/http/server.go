package HTTPServer

import (
	"fmt"
	"github.com/aerosystems/common-service/internal/http_server"
	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"
)

type Server struct {
	srv httpserver.Server
}

func NewHTTPServer(
	cfg *httpserver.Config,
	log *logrus.Logger,
) *Server {
	server := &Server{
		port:                   port,
		log:                    log,
		echo:                   echo.New(),
		firebaseAuthMiddleware: middlewares.FirebaseAuthMiddleware,
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
