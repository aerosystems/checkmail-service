// Code generated by Wire. DO NOT EDIT.

//go:generate go run github.com/google/wire/cmd/wire
//go:build !wireinject
// +build !wireinject

package main

import (
	"cloud.google.com/go/firestore"
	"context"
	"firebase.google.com/go/v4/auth"
	"github.com/aerosystems/checkmail-service/internal/adapters"
	"github.com/aerosystems/checkmail-service/internal/common/config"
	"github.com/aerosystems/checkmail-service/internal/common/custom_errors"
	"github.com/aerosystems/checkmail-service/internal/presenters/grpc"
	"github.com/aerosystems/checkmail-service/internal/presenters/http"
	"github.com/aerosystems/checkmail-service/internal/usecases"
	"github.com/aerosystems/checkmail-service/pkg/gcp"
	"github.com/aerosystems/checkmail-service/pkg/gormconn"
	"github.com/aerosystems/checkmail-service/pkg/logger"
	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

// Injectors from wire.go:

//go:generate wire
func InitApp() *App {
	logger := ProvideLogger()
	logrusLogger := ProvideLogrusLogger(logger)
	config := ProvideConfig()
	httpErrorHandler := ProvideErrorHandler(config)
	baseHandler := ProvideBaseHandler(logrusLogger, config)
	entry := ProvideLogrusEntry(logger)
	db := ProvideGormPostgres(entry, config)
	domainRepo := ProvideDomainRepo(db)
	rootDomainRepo := ProvideRootDomainRepo(db)
	domainUsecase := ProvideDomainUsecase(domainRepo, rootDomainRepo)
	domainHandler := ProvideDomainHandler(baseHandler, domainUsecase)
	filterRepo := ProvideFilterRepo(db)
	filterUsecase := ProvideFilterUsecase(rootDomainRepo, filterRepo)
	filterHandler := ProvideFilterHandler(baseHandler, filterUsecase)
	inspectUsecase := ProvideInspectUsecase(logrusLogger, domainRepo, rootDomainRepo, filterRepo)
	checkHandler := ProvideCheckHandler(baseHandler, inspectUsecase)
	reviewRepo := ProvideReviewRepo(db)
	reviewUsecase := ProvideReviewUsecase(reviewRepo, rootDomainRepo)
	reviewHandler := ProvideReviewHandler(baseHandler, reviewUsecase)
	client := ProvideFirestoreClient(config)
	apiAccessRepo := ProvideApiAccessRepo(client)
	cachedApiAccessRepo := ProvideCachedAccessRepo(apiAccessRepo)
	accessUsecase := ProvideAccessUsecase(cachedApiAccessRepo)
	accessHandler := ProvideAccessHandler(accessUsecase)
	handlers := ProvideHTTPServerHandlers(domainHandler, filterHandler, checkHandler, reviewHandler, accessHandler)
	authClient := ProvideFirebaseAuthClient(config)
	firebaseAuth := ProvideFirebaseAuthMiddleware(authClient)
	apiKeyAuth := ProvideApiKeyMiddleware(accessUsecase)
	middlewares := ProvideHTTPServerMiddlewares(firebaseAuth, apiKeyAuth)
	server := ProvideHttpServer(config, logrusLogger, httpErrorHandler, handlers, middlewares)
	grpcServerCheckHandler := ProvideGRPCCheckHandler(inspectUsecase)
	grpcServerServer := ProvideGRPCServer(logrusLogger, config, grpcServerCheckHandler)
	app := ProvideApp(logrusLogger, config, server, grpcServerServer)
	return app
}

func ProvideApp(log *logrus.Logger, cfg *config.Config, httpServer *HTTPServer.Server, grpcServer *GRPCServer.Server) *App {
	app := NewApp(log, cfg, httpServer, grpcServer)
	return app
}

func ProvideLogger() *logger.Logger {
	loggerLogger := logger.NewLogger()
	return loggerLogger
}

func ProvideConfig() *config.Config {
	configConfig := config.NewConfig()
	return configConfig
}

func ProvideDomainHandler(baseHandler *HTTPServer.BaseHandler, domainUsecase HTTPServer.DomainUsecase) *HTTPServer.DomainHandler {
	domainHandler := HTTPServer.NewDomainHandler(baseHandler, domainUsecase)
	return domainHandler
}

func ProvideFilterHandler(baseHandler *HTTPServer.BaseHandler, filterUsecase HTTPServer.FilterUsecase) *HTTPServer.FilterHandler {
	filterHandler := HTTPServer.NewFilterHandler(baseHandler, filterUsecase)
	return filterHandler
}

func ProvideCheckHandler(baseHandler *HTTPServer.BaseHandler, inspectUsecase HTTPServer.InspectUsecase) *HTTPServer.CheckHandler {
	checkHandler := HTTPServer.NewCheckHandler(baseHandler, inspectUsecase)
	return checkHandler
}

func ProvideReviewHandler(baseHandler *HTTPServer.BaseHandler, reviewUsecase HTTPServer.ReviewUsecase) *HTTPServer.ReviewHandler {
	reviewHandler := HTTPServer.NewReviewHandler(baseHandler, reviewUsecase)
	return reviewHandler
}

func ProvideDomainUsecase(domainRepo usecases.DomainRepository, rootDomainRepo usecases.RootDomainRepository) *usecases.DomainUsecase {
	domainUsecase := usecases.NewDomainUsecase(domainRepo, rootDomainRepo)
	return domainUsecase
}

func ProvideFilterUsecase(rootDomainRepo usecases.RootDomainRepository, filterRepo usecases.FilterRepository) *usecases.FilterUsecase {
	filterUsecase := usecases.NewFilterUsecase(rootDomainRepo, filterRepo)
	return filterUsecase
}

func ProvideInspectUsecase(log *logrus.Logger, domainRepo usecases.DomainRepository, rootDomainRepo usecases.RootDomainRepository, filterRepo usecases.FilterRepository) *usecases.InspectUsecase {
	inspectUsecase := usecases.NewInspectUsecase(log, domainRepo, rootDomainRepo, filterRepo)
	return inspectUsecase
}

func ProvideReviewUsecase(domainReviewRepo usecases.ReviewRepository, rootDomainRepo usecases.RootDomainRepository) *usecases.ReviewUsecase {
	reviewUsecase := usecases.NewReviewUsecase(domainReviewRepo, rootDomainRepo)
	return reviewUsecase
}

func ProvideDomainRepo(db *gorm.DB) *adapters.DomainRepo {
	domainRepo := adapters.NewDomainRepo(db)
	return domainRepo
}

func ProvideRootDomainRepo(db *gorm.DB) *adapters.RootDomainRepo {
	rootDomainRepo := adapters.NewRootDomainRepo(db)
	return rootDomainRepo
}

func ProvideFilterRepo(db *gorm.DB) *adapters.FilterRepo {
	filterRepo := adapters.NewFilterRepo(db)
	return filterRepo
}

func ProvideReviewRepo(db *gorm.DB) *adapters.ReviewRepo {
	reviewRepo := adapters.NewReviewRepo(db)
	return reviewRepo
}

func ProvideApiAccessRepo(client *firestore.Client) *adapters.ApiAccessRepo {
	apiAccessRepo := adapters.NewApiAccessRepo(client)
	return apiAccessRepo
}

func ProvideCachedAccessRepo(apiAccessRepo *adapters.ApiAccessRepo) *adapters.CachedApiAccessRepo {
	cachedApiAccessRepo := adapters.NewCachedApiAccessRepo(apiAccessRepo)
	return cachedApiAccessRepo
}

func ProvideAccessUsecase(apiAccessRepo usecases.ApiAccessRepository) *usecases.AccessUsecase {
	accessUsecase := usecases.NewAccessUsecase(apiAccessRepo)
	return accessUsecase
}

func ProvideApiKeyMiddleware(accessUsecase HTTPServer.AccessUsecase) *HTTPServer.ApiKeyAuth {
	apiKeyAuth := HTTPServer.NewApiKeyAuth(accessUsecase)
	return apiKeyAuth
}

func ProvideAccessHandler(accessUsecase HTTPServer.AccessUsecase) *HTTPServer.AccessHandler {
	accessHandler := HTTPServer.NewAccessHandler(accessUsecase)
	return accessHandler
}

func ProvideGRPCCheckHandler(inspectUsecase GRPCServer.InspectUsecase) *GRPCServer.CheckHandler {
	checkHandler := GRPCServer.NewCheckHandler(inspectUsecase)
	return checkHandler
}

// wire.go:

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

func ProvideFirestoreClient(cfg *config.Config) *firestore.Client {
	ctx := context.Background()
	client, err := firestore.NewClient(ctx, cfg.GcpProjectId)
	if err != nil {
		panic(err)
	}
	return client
}

func ProvideFirebaseAuthClient(cfg *config.Config) *auth.Client {
	app, err := gcp.NewFirebaseApp(cfg.GcpProjectId, cfg.GoogleApplicationCredentials)
	if err != nil {
		panic(err)
	}
	return app.Client
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

func ProvideGRPCServer(log *logrus.Logger, cfg *config.Config, checkHandler *GRPCServer.CheckHandler) *GRPCServer.Server {
	return GRPCServer.NewGRPCServer(cfg.Port, log, checkHandler)
}
