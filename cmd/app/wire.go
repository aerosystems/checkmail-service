//go:build wireinject
// +build wireinject

package main

import (
	"github.com/aerosystems/checkmail-service/internal/config"
	HTTPServer "github.com/aerosystems/checkmail-service/internal/http"
	"github.com/aerosystems/checkmail-service/internal/presenters/rest"
	RPCServer "github.com/aerosystems/checkmail-service/internal/presenters/rpc"
	"github.com/aerosystems/checkmail-service/internal/repository/pg"
	rpcRepo "github.com/aerosystems/checkmail-service/internal/repository/rpc"
	"github.com/aerosystems/checkmail-service/internal/usecases"
	GormPostgres "github.com/aerosystems/checkmail-service/pkg/gorm_postgres"
	"github.com/aerosystems/checkmail-service/pkg/logger"
	OAuthService "github.com/aerosystems/checkmail-service/pkg/oauth_service"
	"github.com/google/wire"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

//go:generate wire
func InitializeApp() *App {
	panic(wire.Build(
		wire.Bind(new(rest.DomainUsecase), new(*usecases.DomainUsecase)),
		wire.Bind(new(rest.FilterUsecase), new(*usecases.FilterUsecase)),
		wire.Bind(new(rest.InspectUsecase), new(*usecases.InspectUsecase)),
		wire.Bind(new(RPCServer.InspectUsecase), new(*usecases.InspectUsecase)),
		wire.Bind(new(rest.ReviewUsecase), new(*usecases.ReviewUsecase)),
		wire.Bind(new(usecases.DomainRepository), new(*pg.DomainRepo)),
		wire.Bind(new(usecases.RootDomainRepository), new(*pg.RootDomainRepo)),
		wire.Bind(new(usecases.FilterRepository), new(*pg.FilterRepo)),
		wire.Bind(new(usecases.ReviewRepository), new(*pg.ReviewRepo)),
		wire.Bind(new(usecases.ProjectRepository), new(*rpcRepo.ProjectRepo)),
		wire.Bind(new(HTTPServer.TokenService), new(*OAuthService.AccessTokenService)),
		ProvideApp,
		ProvideLogger,
		ProvideConfig,
		ProvideHTTPServer,
		ProvideRPCServer,
		ProvideLogrusLogger,
		ProvideLogrusEntry,
		ProvideGormPostgres,
		ProvideBaseHandler,
		ProvideDomainHandler,
		ProvideFilterHandler,
		ProvideInspectHandler,
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

func ProvideApp(log *logrus.Logger, cfg *config.Config, httpServer *HTTPServer.Server, rpcServer *RPCServer.Server) *App {
	panic(wire.Build(NewApp))
}

func ProvideLogger() *logger.Logger {
	panic(wire.Build(logger.NewLogger))
}

func ProvideConfig() *config.Config {
	panic(wire.Build(config.NewConfig))
}

func ProvideHTTPServer(log *logrus.Logger, cfg *config.Config, domainHandler *rest.DomainHandler, filterHandler *rest.FilterHandler, inspectHandler *rest.InspectHandler, reviewHandler *rest.ReviewHandler, tokenService HTTPServer.TokenService) *HTTPServer.Server {
	panic(wire.Build(HTTPServer.NewServer))
}

func ProvideRPCServer(log *logrus.Logger, inspectUsecase RPCServer.InspectUsecase) *RPCServer.Server {
	panic(wire.Build(RPCServer.NewServer))
}

func ProvideLogrusEntry(log *logger.Logger) *logrus.Entry {
	return logrus.NewEntry(log.Logger)
}

func ProvideLogrusLogger(log *logger.Logger) *logrus.Logger {
	return log.Logger
}

func ProvideGormPostgres(e *logrus.Entry, cfg *config.Config) *gorm.DB {
	return GormPostgres.NewClient(e, cfg.PostgresDSN)
}

func ProvideBaseHandler(log *logrus.Logger, cfg *config.Config) *rest.BaseHandler {
	return rest.NewBaseHandler(log, cfg.Mode)
}

func ProvideDomainHandler(baseHandler *rest.BaseHandler, domainUsecase rest.DomainUsecase) *rest.DomainHandler {
	panic(wire.Build(rest.NewDomainHandler))
}

func ProvideFilterHandler(baseHandler *rest.BaseHandler, filterUsecase rest.FilterUsecase) *rest.FilterHandler {
	panic(wire.Build(rest.NewFilterHandler))
}

func ProvideInspectHandler(baseHandler *rest.BaseHandler, inspectUsecase rest.InspectUsecase) *rest.InspectHandler {
	panic(wire.Build(rest.NewInspectHandler))
}

func ProvideReviewHandler(baseHandler *rest.BaseHandler, reviewUsecase rest.ReviewUsecase) *rest.ReviewHandler {
	panic(wire.Build(rest.NewReviewHandler))
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

func ProvideProjectRepo() *rpcRepo.ProjectRepo {
	panic(wire.Build(rpcRepo.NewProjectRepo))
}

func ProvideAccessTokenService(cfg *config.Config) *OAuthService.AccessTokenService {
	return OAuthService.NewAccessTokenService(cfg.AccessSecret)
}
