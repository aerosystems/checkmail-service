//go:build wireinject
// +build wireinject

package main

import (
	"cloud.google.com/go/firestore"
	"context"
	"firebase.google.com/go/auth"
	"github.com/aerosystems/checkmail-service/internal/config"
	"github.com/aerosystems/checkmail-service/internal/infrastructure/repository/fire"
	"github.com/aerosystems/checkmail-service/internal/infrastructure/repository/pg"
	"github.com/aerosystems/checkmail-service/internal/models"
	HttpServer "github.com/aerosystems/checkmail-service/internal/presenters/http"
	"github.com/aerosystems/checkmail-service/internal/presenters/http/handlers"
	"github.com/aerosystems/checkmail-service/internal/presenters/http/handlers/check"
	"github.com/aerosystems/checkmail-service/internal/presenters/http/handlers/domain"
	"github.com/aerosystems/checkmail-service/internal/presenters/http/handlers/filter"
	"github.com/aerosystems/checkmail-service/internal/presenters/http/handlers/review"
	"github.com/aerosystems/checkmail-service/internal/presenters/http/middleware"
	RpcServer "github.com/aerosystems/checkmail-service/internal/presenters/rpc"
	"github.com/aerosystems/checkmail-service/internal/usecases"
	"github.com/aerosystems/checkmail-service/pkg/firebase"
	GormPostgres "github.com/aerosystems/checkmail-service/pkg/gorm_postgres"
	"github.com/aerosystems/checkmail-service/pkg/logger"
	"github.com/google/wire"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

//go:generate wire
func InitApp() *App {
	panic(wire.Build(
		wire.Bind(new(middleware.AccessUsecase), new(*usecases.AccessUsecase)),
		wire.Bind(new(handlers.DomainUsecase), new(*usecases.DomainUsecase)),
		wire.Bind(new(handlers.FilterUsecase), new(*usecases.FilterUsecase)),
		wire.Bind(new(handlers.InspectUsecase), new(*usecases.InspectUsecase)),
		wire.Bind(new(RpcServer.InspectUsecase), new(*usecases.InspectUsecase)),
		wire.Bind(new(handlers.ReviewUsecase), new(*usecases.ReviewUsecase)),
		wire.Bind(new(usecases.DomainRepository), new(*pg.DomainRepo)),
		wire.Bind(new(usecases.RootDomainRepository), new(*pg.RootDomainRepo)),
		wire.Bind(new(usecases.FilterRepository), new(*pg.FilterRepo)),
		wire.Bind(new(usecases.ReviewRepository), new(*pg.ReviewRepo)),
		wire.Bind(new(usecases.ApiAccessRepository), new(*fire.ApiAccessRepo)),
		ProvideApp,
		ProvideLogger,
		ProvideConfig,
		ProvideHttpServer,
		ProvideRpcServer,
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
	))
}

func ProvideApp(log *logrus.Logger, cfg *config.Config, httpServer *HttpServer.Server, rpcServer *RpcServer.Server) *App {
	panic(wire.Build(NewApp))
}

func ProvideLogger() *logger.Logger {
	panic(wire.Build(logger.NewLogger))
}

func ProvideConfig() *config.Config {
	panic(wire.Build(config.NewConfig))
}

func ProvideHttpServer(log *logrus.Logger, cfg *config.Config, firebaseAuthMiddleware *middleware.FirebaseAuth, apiKeyAuthMiddleware *middleware.ApiKeyAuth, domainHandler *domain.Handler, filterHandler *filter.Handler, checkHandler *check.Handler, reviewHandler *review.Handler) *HttpServer.Server {
	panic(wire.Build(HttpServer.NewServer))
}

func ProvideRpcServer(log *logrus.Logger, inspectUsecase RpcServer.InspectUsecase) *RpcServer.Server {
	panic(wire.Build(RpcServer.NewServer))
}

func ProvideLogrusEntry(log *logger.Logger) *logrus.Entry {
	return logrus.NewEntry(log.Logger)
}

func ProvideLogrusLogger(log *logger.Logger) *logrus.Logger {
	return log.Logger
}

func ProvideGormPostgres(e *logrus.Entry, cfg *config.Config) *gorm.DB {
	db := GormPostgres.NewClient(e, cfg.PostgresDSN)
	if err := db.AutoMigrate(&models.Domain{}, &models.RootDomain{}, &models.Filter{}, &models.Review{}); err != nil { // TODO: Move to migration
		panic(err)
	}
	return db
}

func ProvideBaseHandler(log *logrus.Logger, cfg *config.Config) *handlers.BaseHandler {
	return handlers.NewBaseHandler(log, cfg.Mode)
}

func ProvideDomainHandler(baseHandler *handlers.BaseHandler, domainUsecase handlers.DomainUsecase) *domain.Handler {
	panic(wire.Build(domain.NewHandler))
}

func ProvideFilterHandler(baseHandler *handlers.BaseHandler, filterUsecase handlers.FilterUsecase) *filter.Handler {
	panic(wire.Build(filter.NewHandler))
}

func ProvideCheckHandler(baseHandler *handlers.BaseHandler, inspectUsecase handlers.InspectUsecase) *check.Handler {
	panic(wire.Build(check.NewHandler))
}

func ProvideReviewHandler(baseHandler *handlers.BaseHandler, reviewUsecase handlers.ReviewUsecase) *review.Handler {
	panic(wire.Build(review.NewHandler))
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

func ProvideDomainRepo(db *gorm.DB) *pg.DomainRepo {
	panic(wire.Build(pg.NewDomainRepo))
}

func ProvideRootDomainRepo(db *gorm.DB) *pg.RootDomainRepo {
	panic(wire.Build(pg.NewRootDomainRepo))
}

func ProvideFilterRepo(db *gorm.DB) *pg.FilterRepo {
	panic(wire.Build(pg.NewFilterRepo))
}

func ProvideReviewRepo(db *gorm.DB) *pg.ReviewRepo {
	panic(wire.Build(pg.NewReviewRepo))
}

func ProvideFirestoreClient(cfg *config.Config) *firestore.Client {
	ctx := context.Background()
	client, err := firestore.NewClient(ctx, cfg.GcpProjectId)
	if err != nil {
		panic(err)
	}
	return client
}

func ProvideApiAccessRepo(client *firestore.Client) *fire.ApiAccessRepo {
	panic(wire.Build(fire.NewApiAccessRepo))
}

func ProvideAccessUsecase(apiAccessRepo usecases.ApiAccessRepository) *usecases.AccessUsecase {
	panic(wire.Build(usecases.NewAccessUsecase))
}

func ProvideApiKeyMiddleware(accessUsecase middleware.AccessUsecase) *middleware.ApiKeyAuth {
	panic(wire.Build(middleware.NewApiKeyAuth))
}

func ProvideFirebaseAuthClient(cfg *config.Config) *auth.Client {
	app, err := firebaseApp.NewApp(cfg.GcpProjectId, cfg.GoogleApplicationCredentials)
	if err != nil {
		panic(err)
	}
	return app.Client
}

func ProvideFirebaseAuthMiddleware(client *auth.Client) *middleware.FirebaseAuth {
	return middleware.NewFirebaseAuth(client)
}
