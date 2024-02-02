package main

import (
	"github.com/aerosystems/checkmail-service/internal/middleware"
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
) *App {
	wire.Build(
		rest.NewBaseHandler,
		rest.NewDomainHandler,
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
		OAuthService.NewAccessTokenService,
		middleware.NewOAuthMiddleware,
		middleware.NewBasicAuthMiddleware,
	)
	return &App{}
}
