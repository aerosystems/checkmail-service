//go:build wireinject
// +build wireinject

package main

import (
	"github.com/aerosystems/checkmail-service/internal/config"
	"github.com/aerosystems/checkmail-service/internal/infrastructure/repository/pg"
	rpcRepo "github.com/aerosystems/checkmail-service/internal/infrastructure/repository/rpc"
	"github.com/aerosystems/checkmail-service/internal/models"
	HttpServer "github.com/aerosystems/checkmail-service/internal/presenters/http"
	"github.com/aerosystems/checkmail-service/internal/presenters/http/handlers"
	"github.com/aerosystems/checkmail-service/internal/presenters/http/handlers/check"
	"github.com/aerosystems/checkmail-service/internal/presenters/http/handlers/domain"
	"github.com/aerosystems/checkmail-service/internal/presenters/http/handlers/filter"
	"github.com/aerosystems/checkmail-service/internal/presenters/http/handlers/review"
	RpcServer "github.com/aerosystems/checkmail-service/internal/presenters/rpc"
	"github.com/aerosystems/checkmail-service/internal/usecases"
	GormPostgres "github.com/aerosystems/checkmail-service/pkg/gorm_postgres"
	"github.com/aerosystems/checkmail-service/pkg/logger"
	OAuthService "github.com/aerosystems/checkmail-service/pkg/oauth"
	RpcClient "github.com/aerosystems/checkmail-service/pkg/rpc_client"
	"github.com/google/wire"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

//go:generate wire
func InitApp() *App {
	panic(wire.Build(
		wire.Bind(new(handlers.DomainUsecase), new(*usecases.DomainUsecase)),
		wire.Bind(new(handlers.FilterUsecase), new(*usecases.FilterUsecase)),
		wire.Bind(new(handlers.InspectUsecase), new(*usecases.InspectUsecase)),
		wire.Bind(new(RpcServer.InspectUsecase), new(*usecases.InspectUsecase)),
		wire.Bind(new(handlers.ReviewUsecase), new(*usecases.ReviewUsecase)),
		wire.Bind(new(usecases.DomainRepository), new(*pg.DomainRepo)),
		wire.Bind(new(usecases.RootDomainRepository), new(*pg.RootDomainRepo)),
		wire.Bind(new(usecases.FilterRepository), new(*pg.FilterRepo)),
		wire.Bind(new(usecases.ReviewRepository), new(*pg.ReviewRepo)),
		wire.Bind(new(usecases.ProjectRepository), new(*rpcRepo.ProjectRepo)),
		wire.Bind(new(HttpServer.TokenService), new(*OAuthService.AccessTokenService)),
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
		ProvideProjectRepo,
		ProvideAccessTokenService,
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

func ProvideHttpServer(log *logrus.Logger, cfg *config.Config, domainHandler *domain.Handler, filterHandler *filter.Handler, checkHandler *check.Handler, reviewHandler *review.Handler, tokenService HttpServer.TokenService) *HttpServer.Server {
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

func ProvideFilterUsecase(rootDomainRepo usecases.RootDomainRepository, filterRepo usecases.FilterRepository, projectRepo usecases.ProjectRepository) *usecases.FilterUsecase {
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

func ProvideProjectRepo(cfg *config.Config) *rpcRepo.ProjectRepo {
	rpcClient := RpcClient.NewClient("tcp", cfg.ProjectServiceRpcAddress)
	return rpcRepo.NewProjectRepo(rpcClient)
}

func ProvideAccessTokenService(cfg *config.Config) *OAuthService.AccessTokenService {
	return OAuthService.NewAccessTokenService(cfg.AccessSecret)
}
