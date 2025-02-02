//go:build wireinject
// +build wireinject

package main

import (
	"github.com/aerosystems/checkmail-service/internal/adapters"
	GRPCServer "github.com/aerosystems/checkmail-service/internal/ports/grpc"
	HTTPServer "github.com/aerosystems/checkmail-service/internal/ports/http"
	"github.com/aerosystems/checkmail-service/internal/usecases"

	"github.com/aerosystems/common-service/logger"
	"github.com/aerosystems/common-service/pkg/gcpclient"
	"github.com/aerosystems/common-service/pkg/gormclient"
	"github.com/aerosystems/common-service/presenters/grpcserver"
	"github.com/aerosystems/common-service/presenters/httpserver"

	"firebase.google.com/go/v4/auth"
	"github.com/google/wire"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

//go:generate wire
func InitApp() *App {
	panic(wire.Build(
		wire.Bind(new(GRPCServer.InspectUsecase), new(*usecases.InspectUsecase)),
		wire.Bind(new(HTTPServer.AccessUsecase), new(*usecases.AccessUsecase)),
		wire.Bind(new(HTTPServer.ManageUsecase), new(*usecases.ManageUsecase)),
		wire.Bind(new(HTTPServer.InspectUsecase), new(*usecases.InspectUsecase)),
		wire.Bind(new(usecases.DomainRepository), new(*adapters.DomainRepo)),
		wire.Bind(new(usecases.FilterRepository), new(*adapters.FilterRepo)),
		wire.Bind(new(usecases.AccessRepository), new(*adapters.AccessRepo)),
		ProvideApp,
		ProvideLogger,
		ProvideConfig,
		ProvideHTTPServer,
		ProvideLogrusLogger,
		ProvideGORMPostgres,
		ProvideBaseHandler,
		ProvideDomainHandler,
		ProvideFilterHandler,
		ProvideCheckHandler,
		ProvideManageUsecase,
		ProvideInspectUsecase,
		ProvideDomainRepo,
		ProvideFilterRepo,
		ProvideAccessUsecase,
		ProvideAccessRepo,
		ProvideFirebaseAuthClient,
		ProvideFirebaseAuthMiddleware,
		ProvideAccessHandler,
		ProvideGRPCCheckHandler,
		ProvideGRPCServer,
	))
}

func ProvideApp(log *logrus.Logger, cfg *Config, httpServer *HTTPServer.Server, grpcServer *GRPCServer.Server) *App {
	panic(wire.Build(NewApp))
}

func ProvideLogger() *logger.Logger {
	panic(wire.Build(logger.NewLogger))
}

func ProvideConfig() *Config {
	panic(wire.Build(NewConfig))
}

func ProvideHTTPServer(cfg *Config, log *logrus.Logger, firebaseAuth *HTTPServer.FirebaseAuth,
	checkHandler *HTTPServer.CheckHandler, accessHandler *HTTPServer.AccessHandler,
	domainHandler *HTTPServer.DomainHandler, filterHandler *HTTPServer.FilterHandler) *HTTPServer.Server {
	return HTTPServer.NewHTTPServer(&HTTPServer.Config{
		Config: httpserver.Config{
			Host: cfg.Host,
			Port: cfg.Port,
		},
		Mode: cfg.Mode,
	}, log, firebaseAuth, checkHandler, accessHandler, domainHandler, filterHandler)
}

func ProvideLogrusLogger(log *logger.Logger) *logrus.Logger {
	return log.Logger
}

func ProvideGORMPostgres(log *logrus.Logger, cfg *Config) *gorm.DB {
	db := gormclient.NewPostgresDB(log, cfg.PostgresDSN)
	if err := adapters.AutoMigrateGORM(db); err != nil {
		panic(err)
	}
	return db
}

func ProvideBaseHandler(log *logrus.Logger, cfg *Config) *HTTPServer.BaseHandler {
	return HTTPServer.NewBaseHandler(log, cfg.Mode)
}

func ProvideDomainHandler(baseHandler *HTTPServer.BaseHandler, domainUsecase HTTPServer.ManageUsecase) *HTTPServer.DomainHandler {
	panic(wire.Build(HTTPServer.NewDomainHandler))
}

func ProvideFilterHandler(baseHandler *HTTPServer.BaseHandler, manageUsecase HTTPServer.ManageUsecase) *HTTPServer.FilterHandler {
	panic(wire.Build(HTTPServer.NewFilterHandler))
}

func ProvideCheckHandler(baseHandler *HTTPServer.BaseHandler, inspectUsecase HTTPServer.InspectUsecase) *HTTPServer.CheckHandler {
	panic(wire.Build(HTTPServer.NewCheckHandler))
}

func ProvideReviewHandler(baseHandler *HTTPServer.BaseHandler, reviewUsecase HTTPServer.ReviewUsecase) *HTTPServer.ReviewHandler {
	panic(wire.Build(HTTPServer.NewReviewHandler))
}

func ProvideManageUsecase(domainRepo usecases.DomainRepository, filterRepo usecases.FilterRepository) *usecases.ManageUsecase {
	panic(wire.Build(usecases.NewManageUsecase))
}

func ProvideInspectUsecase(log *logrus.Logger, accessRepo usecases.AccessRepository, domainRepo usecases.DomainRepository, filterRepo usecases.FilterRepository) *usecases.InspectUsecase {
	panic(wire.Build(usecases.NewInspectUsecase))
}

func ProvideReviewUsecase(domainReviewRepo usecases.ReviewRepository) *usecases.ReviewUsecase {
	panic(wire.Build(usecases.NewReviewUsecase))
}

func ProvideDomainRepo(db *gorm.DB) *adapters.DomainRepo {
	panic(wire.Build(adapters.NewDomainRepo))
}

func ProvideFilterRepo(db *gorm.DB) *adapters.FilterRepo {
	panic(wire.Build(adapters.NewFilterRepo))
}

func ProvideReviewRepo(db *gorm.DB) *adapters.ReviewRepo {
	panic(wire.Build(adapters.NewReviewRepo))
}

func ProvideAccessRepo(db *gorm.DB) *adapters.AccessRepo {
	panic(wire.Build(adapters.NewAccessRepo))
}

func ProvideAccessUsecase(apiAccessRepo usecases.AccessRepository) *usecases.AccessUsecase {
	panic(wire.Build(usecases.NewAccessUsecase))
}

func ProvideFirebaseAuthClient(cfg *Config) *auth.Client {
	client, err := gcpclient.NewFirebaseClient(cfg.GcpProjectId, cfg.GoogleApplicationCredentials)
	if err != nil {
		panic(err)
	}
	return client
}

func ProvideFirebaseAuthMiddleware(client *auth.Client) *HTTPServer.FirebaseAuth {
	return HTTPServer.NewFirebaseAuth(client)
}

func ProvideAccessHandler(accessUsecase HTTPServer.AccessUsecase) *HTTPServer.AccessHandler {
	panic(wire.Build(HTTPServer.NewAccessHandler))
}

func ProvideGRPCCheckHandler(inspectUsecase GRPCServer.InspectUsecase) *GRPCServer.CheckHandler {
	panic(wire.Build(GRPCServer.NewCheckHandler))
}

func ProvideGRPCServer(log *logrus.Logger, cfg *Config, checkHandler *GRPCServer.CheckHandler) *GRPCServer.Server {
	return GRPCServer.NewGRPCServer(&grpcserver.Config{Host: cfg.Host, Port: cfg.Port}, log, checkHandler)
}
