package HttpServer

import (
	"fmt"
	"github.com/aerosystems/checkmail-service/internal/presenters/http/handlers/check"
	"github.com/aerosystems/checkmail-service/internal/presenters/http/handlers/domain"
	"github.com/aerosystems/checkmail-service/internal/presenters/http/handlers/filter"
	"github.com/aerosystems/checkmail-service/internal/presenters/http/handlers/review"
	"github.com/aerosystems/checkmail-service/internal/presenters/http/middleware"
	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"
)

const webPort = 80

type Server struct {
	log                    *logrus.Logger
	echo                   *echo.Echo
	firebaseAuthMiddleware *middleware.FirebaseAuth
	xApiKeyAuthMiddleware  *middleware.ApiKeyAuth
	domainHandler          *domain.Handler
	filterHandler          *filter.Handler
	checkHandler           *check.Handler
	reviewHandler          *review.Handler
	tokenService           TokenService
}

func NewServer(
	log *logrus.Logger,
	firebaseAuthMiddleware *middleware.FirebaseAuth,
	xApiKeyAuthMiddleware *middleware.ApiKeyAuth,
	domainHandler *domain.Handler,
	filterHandler *filter.Handler,
	checkHandler *check.Handler,
	reviewHandler *review.Handler,
	tokenService TokenService,
) *Server {
	return &Server{
		log:                    log,
		echo:                   echo.New(),
		firebaseAuthMiddleware: firebaseAuthMiddleware,
		xApiKeyAuthMiddleware:  xApiKeyAuthMiddleware,
		domainHandler:          domainHandler,
		filterHandler:          filterHandler,
		checkHandler:           checkHandler,
		reviewHandler:          reviewHandler,
		tokenService:           tokenService,
	}
}

func (s *Server) Run() error {
	s.setupMiddleware()
	s.setupRoutes()
	s.log.Infof("starting HTTP server checkmail-service on port %d\n", webPort)
	return s.echo.Start(fmt.Sprintf(":%d", webPort))
}
