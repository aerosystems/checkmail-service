package main

import (
	"fmt"
	"github.com/aerosystems/checkmail-service/internal/middleware"
	"github.com/aerosystems/checkmail-service/internal/models"
	"github.com/aerosystems/checkmail-service/internal/presenters/rest"
	"github.com/aerosystems/checkmail-service/internal/repository"
	RPCServer "github.com/aerosystems/checkmail-service/internal/rpc_server"
	"github.com/aerosystems/checkmail-service/internal/usecases"
	GormPostgres "github.com/aerosystems/checkmail-service/pkg/gorm_postgres"
	"github.com/aerosystems/checkmail-service/pkg/logger"
	OAuthService "github.com/aerosystems/checkmail-service/pkg/oauth_service"
	"github.com/sirupsen/logrus"
	"net/rpc"
	"os"
)

const (
	webPort = 80
	rpcPort = 5001
)

// @title Checkmail Service
// @version 1.0.7
// @description A part of microservice infrastructure, who responsible for store and check email domains in black/whitelists

// @contact.name Artem Kostenko
// @contact.url https://github.com/aerosystems

// @license.name Apache 2.0
// @license.url https://www.apache.org/licenses/LICENSE-2.0.html

// @securityDefinitions.apikey X-Api-Key
// @in header
// @name X-Api-Key
// @description Should contain Token, digits and letters, 64 symbols length

// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @description Should contain Access JWT Token, with the Bearer started

// @host gw.verifire.app/checkmail
// @schemes https
// @BasePath /
func main() {
	log := logger.NewLogger(os.Getenv("HOSTNAME"))

	clientGORM := GormPostgres.NewClient(logrus.NewEntry(log.Logger))
	if err := clientGORM.AutoMigrate(&models.Domain{}, &models.RootDomain{}, &models.Filter{}, &models.DomainReview{}); err != nil {
		log.Panic(err)
	}
	domainRepo := repository.NewDomainRepo(clientGORM)
	rootDomainRepo := repository.NewRootDomainRepo(clientGORM)
	filterRepo := repository.NewFilterRepo(clientGORM)
	domainReviewRepo := repository.NewDomainReviewRepo(clientGORM)

	inspectService := usecases.NewInspectService(log.Logger, domainRepo, rootDomainRepo, filterRepo)

	checkmailServer := RPCServer.NewCheckmailServer(rpcPort, inspectService)

	baseHandler := rest.NewBaseHandler(log.Logger, domainRepo, rootDomainRepo, filterRepo, domainReviewRepo, inspectService)

	accessTokenService := OAuthService.NewAccessTokenService(os.Getenv("ACCESS_SECRET"))

	oauthMiddleware := middleware.NewOAuthMiddlewareImpl(accessTokenService)
	basicAuthMiddleware := middleware.NewBasicAuthMiddlewareImpl(os.Getenv("BASIC_AUTH_DOCS_USERNAME"), os.Getenv("BASIC_AUTH_DOCS_PASSWORD"))

	app := NewConfig(baseHandler, oauthMiddleware, basicAuthMiddleware)

	e := app.NewRouter()
	middleware.AddLog(e, log.Logger)

	errChan := make(chan error)

	go func() {
		log.Infof("starting checkmail-service RPC server on port %d\n", rpcPort)
		errChan <- rpc.Register(checkmailServer)
		errChan <- checkmailServer.Listen()
	}()

	go func() {
		log.Infof("starting HTTP server project-service on port %d\n", webPort)
		errChan <- e.Start(fmt.Sprintf(":%d", webPort))
	}()

	for {
		select {
		case err := <-errChan:
			if err != nil {
				log.Fatal(err)
			}
		}
	}
}
