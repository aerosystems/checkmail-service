package HttpServer

import (
	"fmt"
	"github.com/aerosystems/checkmail-service/internal/presenters/http/handlers/access"
	"github.com/aerosystems/checkmail-service/internal/presenters/http/handlers/check"
	"github.com/aerosystems/checkmail-service/internal/presenters/http/handlers/domain"
	"github.com/aerosystems/checkmail-service/internal/presenters/http/handlers/filter"
	"github.com/aerosystems/checkmail-service/internal/presenters/http/handlers/review"
	"github.com/aerosystems/checkmail-service/internal/presenters/http/middleware"
	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"
)

type Server struct {
	echo                   *echo.Echo
	port                   int
	log                    *logrus.Logger
	firebaseAuthMiddleware *middleware.FirebaseAuth
	apiKeyAuthMiddleware   *middleware.ApiKeyAuth
	domainHandler          *domain.Handler
	filterHandler          *filter.Handler
	checkHandler           *check.Handler
	reviewHandler          *review.Handler
	accessHandler          *access.Handler
}

type Handlers struct {
	DomainHandler *domain.Handler
	FilterHandler *filter.Handler
	CheckHandler  *check.Handler
	ReviewHandler *review.Handler
	AccessHandler *access.Handler
}

type Middlewares struct {
	FirebaseAuthMiddleware *middleware.FirebaseAuth
	ApiKeyAuthMiddleware   *middleware.ApiKeyAuth
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
