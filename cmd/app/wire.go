//go:build wireinject
// +build wireinject

package main

import (
	"cloud.google.com/go/firestore"
	"context"
	"firebase.google.com/go/v4/auth"
	"github.com/aerosystems/checkmail-service/internal/adapters"
	"github.com/aerosystems/checkmail-service/internal/common/config"
	CustomErrors "github.com/aerosystems/checkmail-service/internal/common/custom_errors"
	GRPCServer "github.com/aerosystems/checkmail-service/internal/presenters/grpc"
	HttpServer "github.com/aerosystems/checkmail-service/internal/presenters/http"
	"github.com/aerosystems/checkmail-service/internal/usecases"
	"github.com/aerosystems/checkmail-service/pkg/gormconn"
	"github.com/aerosystems/checkmail-service/pkg/logger"
	"github.com/aerosystems/project-service/pkg/gcp"
	"github.com/google/wire"
	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

//go:generate wire
func InitApp() *App {
	panic(wire.Build(
		wire.Bind(new(GRPCServer.InspectUsecase), new(*usecases.InspectUsecase)),
		wire.Bind(new(HttpServer.AccessUsecase), new(*usecases.AccessUsecase)),
		wire.Bind(new(HttpServer.DomainUsecase), new(*usecases.DomainUsecase)),
		wire.Bind(new(HttpServer.FilterUsecase), new(*usecases.FilterUsecase)),
		wire.Bind(new(HttpServer.InspectUsecase), new(*usecases.InspectUsecase)),
		wire.Bind(new(HttpServer.ReviewUsecase), new(*usecases.ReviewUsecase)),
		wire.Bind(new(usecases.DomainRepository), new(*adapters.DomainRepo)),
		wire.Bind(new(usecases.RootDomainRepository), new(*adapters.RootDomainRepo)),
		wire.Bind(new(usecases.FilterRepository), new(*adapters.FilterRepo)),
		wire.Bind(new(usecases.ReviewRepository), new(*adapters.ReviewRepo)),
		wire.Bind(new(usecases.ApiAccessRepository), new(*adapters.ApiAccessRepo)),
		ProvideApp,
		ProvideLogger,
		ProvideConfig,
		ProvideHttpServer,
		ProvideLogrusLogger,
		ProvideLogrusEntry,
		ProvideGormPostgres,
		ProvideBaseHandler,
		ProvideDomainHandler,
		ProvideFilterHandler,
		ProvideCheckHandler,
		ProvideReviewHandler,
		ProvideDomainUsecase,
		ProvideFilterUsecase,
		ProvideInspectUsecase,
		ProvideReviewUsecase,
		ProvideDomainRepo,
		ProvideRootDomainRepo,
		ProvideFilterRepo,
		ProvideReviewRepo,
		ProvideAccessUsecase,
		ProvideFirestoreClient,
		ProvideApiAccessRepo,
		ProvideApiKeyMiddleware,
		ProvideFirebaseAuthClient,
		ProvideFirebaseAuthMiddleware,
		ProvideErrorHandler,
		ProvideHTTPServerHandlers,
		ProvideHTTPServerMiddlewares,
		ProvideAccessHandler,
		ProvideGRPCCheckHandler,
		ProvideGRPCServer,
	))
}

func ProvideApp(log *logrus.Logger, cfg *config.Config, httpServer *HttpServer.Server, grpcServer *GRPCServer.Server) *App {
	panic(wire.Build(NewApp))
}

func ProvideLogger() *logger.Logger {
	panic(wire.Build(logger.NewLogger))
}

func ProvideConfig() *config.Config {
	panic(wire.Build(config.NewConfig))
}

func ProvideHttpServer(cfg *config.Config, log *logrus.Logger, errorHandler *echo.HTTPErrorHandler, handlers HttpServer.Handlers, middlewares HttpServer.Middlewares) *HttpServer.Server {
	return HttpServer.NewServer(cfg.Port, log, errorHandler, handlers, middlewares)
}

func ProvideLogrusEntry(log *logger.Logger) *logrus.Entry {
	return logrus.NewEntry(log.Logger)
}

func ProvideLogrusLogger(log *logger.Logger) *logrus.Logger {
	return log.Logger
}

func ProvideGormPostgres(e *logrus.Entry, cfg *config.Config) *gorm.DB {
	return gormconn.NewPostgresDB(e, cfg.PostgresDSN)
}

func ProvideBaseHandler(log *logrus.Logger, cfg *config.Config) *HttpServer.BaseHandler {
	return HttpServer.NewBaseHandler(log, cfg.Mode)
}

func ProvideDomainHandler(baseHandler *HttpServer.BaseHandler, domainUsecase HttpServer.DomainUsecase) *HttpServer.DomainHandler {
	panic(wire.Build(HttpServer.NewDomainHandler))
}

func ProvideFilterHandler(baseHandler *HttpServer.BaseHandler, filterUsecase HttpServer.FilterUsecase) *HttpServer.FilterHandler {
	panic(wire.Build(HttpServer.NewFilterHandler))
}

func ProvideCheckHandler(baseHandler *HttpServer.BaseHandler, inspectUsecase HttpServer.InspectUsecase) *HttpServer.CheckHandler {
	panic(wire.Build(HttpServer.NewCheckHandler))
}

func ProvideReviewHandler(baseHandler *HttpServer.BaseHandler, reviewUsecase HttpServer.ReviewUsecase) *HttpServer.ReviewHandler {
	panic(wire.Build(HttpServer.NewReviewHandler))
}

func ProvideDomainUsecase(domainRepo usecases.DomainRepository, rootDomainRepo usecases.RootDomainRepository) *usecases.DomainUsecase {
	panic(wire.Build(usecases.NewDomainUsecase))
}

func ProvideFilterUsecase(rootDomainRepo usecases.RootDomainRepository, filterRepo usecases.FilterRepository) *usecases.FilterUsecase {
	panic(wire.Build(usecases.NewFilterUsecase))
}

func ProvideInspectUsecase(log *logrus.Logger, domainRepo usecases.DomainRepository, rootDomainRepo usecases.RootDomainRepository, filterRepo usecases.FilterRepository) *usecases.InspectUsecase {
	panic(wire.Build(usecases.NewInspectUsecase))
}

func ProvideReviewUsecase(domainReviewRepo usecases.ReviewRepository, rootDomainRepo usecases.RootDomainRepository) *usecases.ReviewUsecase {
	panic(wire.Build(usecases.NewReviewUsecase))
}

func ProvideDomainRepo(db *gorm.DB) *adapters.DomainRepo {
	panic(wire.Build(adapters.NewDomainRepo))
}

func ProvideRootDomainRepo(db *gorm.DB) *adapters.RootDomainRepo {
	panic(wire.Build(adapters.NewRootDomainRepo))
}

func ProvideFilterRepo(db *gorm.DB) *adapters.FilterRepo {
	panic(wire.Build(adapters.NewFilterRepo))
}

func ProvideReviewRepo(db *gorm.DB) *adapters.ReviewRepo {
	panic(wire.Build(adapters.NewReviewRepo))
}

func ProvideFirestoreClient(cfg *config.Config) *firestore.Client {
	ctx := context.Background()
	client, err := firestore.NewClient(ctx, cfg.GcpProjectId)
	if err != nil {
		panic(err)
	}
	return client
}

func ProvideApiAccessRepo(client *firestore.Client) *adapters.ApiAccessRepo {
	panic(wire.Build(adapters.NewApiAccessRepo))
}

func ProvideAccessUsecase(apiAccessRepo usecases.ApiAccessRepository) *usecases.AccessUsecase {
	panic(wire.Build(usecases.NewAccessUsecase))
}

func ProvideApiKeyMiddleware(accessUsecase HttpServer.AccessUsecase) *HttpServer.ApiKeyAuth {
	panic(wire.Build(HttpServer.NewApiKeyAuth))
}

func ProvideFirebaseAuthClient(cfg *config.Config) *auth.Client {
	app, err := gcp.NewFirebaseApp(cfg.GcpProjectId, cfg.GoogleApplicationCredentials)
	if err != nil {
		panic(err)
	}
	return app.Client
}

func ProvideFirebaseAuthMiddleware(client *auth.Client) *HttpServer.FirebaseAuth {
	return HttpServer.NewFirebaseAuth(client)
}

func ProvideErrorHandler(cfg *config.Config) *echo.HTTPErrorHandler {
	errorHandler := CustomErrors.NewEchoErrorHandler(cfg.Mode)
	return &errorHandler
}

func ProvideHTTPServerHandlers(domainHandler *HttpServer.DomainHandler, filterHandler *HttpServer.FilterHandler, checkHandler *HttpServer.CheckHandler, reviewHandler *HttpServer.ReviewHandler, accessHandler *HttpServer.AccessHandler) HttpServer.Handlers {
	return HttpServer.Handlers{
		DomainHandler: domainHandler,
		FilterHandler: filterHandler,
		CheckHandler:  checkHandler,
		ReviewHandler: reviewHandler,
		AccessHandler: accessHandler,
	}
}

func ProvideHTTPServerMiddlewares(firebaseAuthMiddleware *HttpServer.FirebaseAuth, apiKeyAuthMiddleware *HttpServer.ApiKeyAuth) HttpServer.Middlewares {
	return HttpServer.Middlewares{
		FirebaseAuthMiddleware: firebaseAuthMiddleware,
		ApiKeyAuthMiddleware:   apiKeyAuthMiddleware,
	}
}

func ProvideAccessHandler(accessUsecase HttpServer.AccessUsecase) *HttpServer.AccessHandler {
	panic(wire.Build(HttpServer.NewAccessHandler))
}

func ProvideGRPCCheckHandler(inspectUsecase GRPCServer.InspectUsecase) *GRPCServer.CheckHandler {
	panic(wire.Build(GRPCServer.NewCheckHandler))
}

func ProvideGRPCServer(log *logrus.Logger, cfg *config.Config, checkHandler *GRPCServer.CheckHandler) *GRPCServer.Server {
	return GRPCServer.NewGRPCServer(cfg.Port, log, checkHandler)
}
