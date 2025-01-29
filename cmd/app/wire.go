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
	HTTPServer "github.com/aerosystems/checkmail-service/internal/presenters/http"
	"github.com/aerosystems/checkmail-service/internal/usecases"
	"github.com/aerosystems/common-service/pkg/gcp"
	"github.com/aerosystems/common-service/pkg/gormconn"
	"github.com/aerosystems/common-service/pkg/logger"
	"github.com/google/wire"
	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

//go:generate wire
func InitApp() *App {
	panic(wire.Build(
		wire.Bind(new(GRPCServer.InspectUsecase), new(*usecases.InspectUsecase)),
		wire.Bind(new(HTTPServer.AccessUsecase), new(*usecases.AccessUsecase)),
		wire.Bind(new(HTTPServer.DomainUsecase), new(*usecases.DomainUsecase)),
		wire.Bind(new(HTTPServer.FilterUsecase), new(*usecases.FilterUsecase)),
		wire.Bind(new(HTTPServer.InspectUsecase), new(*usecases.InspectUsecase)),
		wire.Bind(new(HTTPServer.ReviewUsecase), new(*usecases.ReviewUsecase)),
		wire.Bind(new(usecases.DomainRepository), new(*adapters.DomainRepo)),
		wire.Bind(new(usecases.FilterRepository), new(*adapters.FilterRepo)),
		wire.Bind(new(usecases.ReviewRepository), new(*adapters.ReviewRepo)),
		wire.Bind(new(usecases.ApiAccessRepository), new(*adapters.CachedApiAccessRepo)),
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
		ProvideFilterRepo,
		ProvideReviewRepo,
		ProvideAccessUsecase,
		ProvideFirestoreClient,
		ProvideApiAccessRepo,
		ProvideCachedAccessRepo,
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

func ProvideApp(log *logrus.Logger, cfg *config.Config, httpServer *HTTPServer.Server, grpcServer *GRPCServer.Server) *App {
	panic(wire.Build(NewApp))
}

func ProvideLogger() *logger.Logger {
	panic(wire.Build(logger.NewLogger))
}

func ProvideConfig() *config.Config {
	panic(wire.Build(config.NewConfig))
}

func ProvideHttpServer(cfg *config.Config, log *logrus.Logger, errorHandler *echo.HTTPErrorHandler, handlers HTTPServer.Handlers, middlewares HTTPServer.Middlewares) *HTTPServer.Server {
	return HTTPServer.NewServer(cfg.Port, log, errorHandler, handlers, middlewares)
}

func ProvideLogrusEntry(log *logger.Logger) *logrus.Entry {
	return logrus.NewEntry(log.Logger)
}

func ProvideLogrusLogger(log *logger.Logger) *logrus.Logger {
	return log.Logger
}

func ProvideGormPostgres(e *logrus.Entry, cfg *config.Config) *gorm.DB {
	db := gormconn.NewPostgresDB(e, cfg.PostgresDSN)
	if err := adapters.AutoMigrateGORM(db); err != nil {
		panic(err)
	}
	return db
}

func ProvideBaseHandler(log *logrus.Logger, cfg *config.Config) *HTTPServer.BaseHandler {
	return HTTPServer.NewBaseHandler(log, cfg.Mode)
}

func ProvideDomainHandler(baseHandler *HTTPServer.BaseHandler, domainUsecase HTTPServer.DomainUsecase) *HTTPServer.DomainHandler {
	panic(wire.Build(HTTPServer.NewDomainHandler))
}

func ProvideFilterHandler(baseHandler *HTTPServer.BaseHandler, filterUsecase HTTPServer.FilterUsecase) *HTTPServer.FilterHandler {
	panic(wire.Build(HTTPServer.NewFilterHandler))
}

func ProvideCheckHandler(baseHandler *HTTPServer.BaseHandler, inspectUsecase HTTPServer.InspectUsecase) *HTTPServer.CheckHandler {
	panic(wire.Build(HTTPServer.NewCheckHandler))
}

func ProvideReviewHandler(baseHandler *HTTPServer.BaseHandler, reviewUsecase HTTPServer.ReviewUsecase) *HTTPServer.ReviewHandler {
	panic(wire.Build(HTTPServer.NewReviewHandler))
}

func ProvideDomainUsecase(domainRepo usecases.DomainRepository) *usecases.DomainUsecase {
	panic(wire.Build(usecases.NewDomainUsecase))
}

func ProvideFilterUsecase(filterRepo usecases.FilterRepository) *usecases.FilterUsecase {
	panic(wire.Build(usecases.NewFilterUsecase))
}

func ProvideInspectUsecase(log *logrus.Logger, domainRepo usecases.DomainRepository, filterRepo usecases.FilterRepository) *usecases.InspectUsecase {
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

func ProvideCachedAccessRepo(apiAccessRepo *adapters.ApiAccessRepo) *adapters.CachedApiAccessRepo {
	panic(wire.Build(adapters.NewCachedApiAccessRepo))
}

func ProvideAccessUsecase(apiAccessRepo usecases.ApiAccessRepository) *usecases.AccessUsecase {
	panic(wire.Build(usecases.NewAccessUsecase))
}

func ProvideApiKeyMiddleware(accessUsecase HTTPServer.AccessUsecase) *HTTPServer.ApiKeyAuth {
	panic(wire.Build(HTTPServer.NewApiKeyAuth))
}

func ProvideFirebaseAuthClient(cfg *config.Config) *auth.Client {
	client, err := gcp.NewFirebaseClient(cfg.GcpProjectId, cfg.GoogleApplicationCredentials)
	if err != nil {
		panic(err)
	}
	return client
}

func ProvideFirebaseAuthMiddleware(client *auth.Client) *HTTPServer.FirebaseAuth {
	return HTTPServer.NewFirebaseAuth(client)
}

func ProvideErrorHandler(cfg *config.Config) *echo.HTTPErrorHandler {
	errorHandler := CustomErrors.NewEchoErrorHandler(cfg.Mode)
	return &errorHandler
}

func ProvideHTTPServerHandlers(domainHandler *HTTPServer.DomainHandler, filterHandler *HTTPServer.FilterHandler, checkHandler *HTTPServer.CheckHandler, reviewHandler *HTTPServer.ReviewHandler, accessHandler *HTTPServer.AccessHandler) HTTPServer.Handlers {
	return HTTPServer.Handlers{
		DomainHandler: domainHandler,
		FilterHandler: filterHandler,
		CheckHandler:  checkHandler,
		ReviewHandler: reviewHandler,
		AccessHandler: accessHandler,
	}
}

func ProvideHTTPServerMiddlewares(firebaseAuthMiddleware *HTTPServer.FirebaseAuth, apiKeyAuthMiddleware *HTTPServer.ApiKeyAuth) HTTPServer.Middlewares {
	return HTTPServer.Middlewares{
		FirebaseAuthMiddleware: firebaseAuthMiddleware,
		ApiKeyAuthMiddleware:   apiKeyAuthMiddleware,
	}
}

func ProvideAccessHandler(accessUsecase HTTPServer.AccessUsecase) *HTTPServer.AccessHandler {
	panic(wire.Build(HTTPServer.NewAccessHandler))
}

func ProvideGRPCCheckHandler(inspectUsecase GRPCServer.InspectUsecase) *GRPCServer.CheckHandler {
	panic(wire.Build(GRPCServer.NewCheckHandler))
}

func ProvideGRPCServer(log *logrus.Logger, cfg *config.Config, checkHandler *GRPCServer.CheckHandler) *GRPCServer.Server {
	return GRPCServer.NewGRPCServer(cfg.Port, log, checkHandler)
}
