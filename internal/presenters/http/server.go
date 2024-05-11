package HttpServer

import (
	"fmt"
	"github.com/aerosystems/checkmail-service/internal/presenters/http/handlers"
	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"
)

const webPort = 80

type Server struct {
	log            *logrus.Logger
	echo           *echo.Echo
	domainHandler  *handlers.DomainHandler
	filterHandler  *handlers.FilterHandler
	inspectHandler *handlers.InspectHandler
	reviewHandler  *handlers.ReviewHandler
	tokenService   TokenService
}

func NewServer(
	log *logrus.Logger,
	domainHandler *handlers.DomainHandler,
	filterHandler *handlers.FilterHandler,
	inspectHandler *handlers.InspectHandler,
	reviewHandler *handlers.ReviewHandler,
	tokenService TokenService,
) *Server {
	return &Server{
		log:            log,
		echo:           echo.New(),
		domainHandler:  domainHandler,
		filterHandler:  filterHandler,
		inspectHandler: inspectHandler,
		reviewHandler:  reviewHandler,
		tokenService:   tokenService,
	}
}

func (s *Server) Run() error {
	s.setupMiddleware()
	s.setupRoutes()
	s.log.Infof("starting HTTP server checkmail-service on port %d\n", webPort)
	return s.echo.Start(fmt.Sprintf(":%d", webPort))
}
