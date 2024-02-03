//go:build wireinject
// +build wireinject

package main

import (
	middleware "github.com/aerosystems/checkmail-service/internal/middleware"
	"github.com/aerosystems/checkmail-service/internal/presenters/rest"
	"github.com/aerosystems/checkmail-service/internal/repository/pg"
	RPCClient "github.com/aerosystems/checkmail-service/internal/repository/rpc"
	"github.com/aerosystems/checkmail-service/internal/usecases"
	OAuthService "github.com/aerosystems/checkmail-service/pkg/oauth_service"
	"github.com/google/wire"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

//go:generate wire
func InitializeApp(
	log *logrus.Logger,
	clientGORM *gorm.DB,
) App {
	wire.Build(
		ProvideApp,
		ProvideBaseHandler,
		ProvideDomainHandler,
		rest.NewFilterHandler,
		rest.NewInspectHandler,
		rest.NewReviewHandler,
		usecases.NewDomainUsecase,
		usecases.NewFilterUsecase,
		usecases.NewInspectUsecase,
		usecases.NewReviewUsecase,
		pg.NewDomainRepo,
		pg.NewRootDomainRepo,
		pg.NewFilterRepo,
		pg.NewReviewRepo,
		RPCClient.NewProjectRepo,
	)
	return App{}
}

func ProvideApp(domainHandler rest.DomainHandler, filterHandler rest.FilterHandler, inspectHandler rest.InspectHandler, reviewHandler rest.ReviewHandler, oauthMiddleware middleware.OAuthMiddleware, basicAuthMiddleware middleware.BasicAuthMiddleware) App {
	panic(wire.Build(wire.Struct(new(App), "*")))
}

func ProvideBaseHandler(log *logrus.Logger) rest.BaseHandler {
	panic(wire.Build(wire.Struct(new(rest.BaseHandler), "*")))
}

func ProvideDomainHandler(baseHandler rest.BaseHandler, domainUsecase usecases.DomainUsecase) rest.DomainHandler {
	panic(wire.Build(wire.Struct(new(rest.DomainHandler), "*")))
}

func ProvideOAuthMiddleware(accessTokenService OAuthService.AccessTokenService) middleware.OAuthMiddleware {
	panic(wire.Build(wire.Struct(new(middleware.OAuthMiddleware), "*")))
}

func ProvideBasicAuthMiddleware(username string, password string) middleware.BasicAuthMiddleware {
	panic(wire.Build(wire.Struct(new(middleware.BasicAuthMiddleware), "*")))
}

func ProvideAccessTokenService() OAuthService.AccessTokenService {
	panic(wire.Build(wire.Struct(new(OAuthService.AccessTokenService), "*")))
}
