package main

import (
	"fmt"
	"github.com/labstack/gommon/log"
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

// @host gw.verifire.dev/checkmail
// @schemes https
// @BasePath /
func main() {

	//if err := clientGORM.AutoMigrate(&models.Domain{}, &models.RootDomain{}, &models.Filter{}, &models.Review{}); err != nil {
	//	log.Panic(err)
	//}

	//domainRepo := pg.NewDomainRepo(clientGORM)
	//rootDomainRepo := pg.NewRootDomainRepo(clientGORM)
	//filterRepo := pg.NewFilterRepo(clientGORM)
	//domainReviewRepo := pg.NewReviewRepo(clientGORM)
	//
	//domainUsecase := usecases.NewDomainUsecase(domainRepo, rootDomainRepo)
	//filterUsecase := usecases.NewFilterUsecase(filterRepo)
	//inspectUsecase := usecases.NewInspectUsecase(log.Logger, domainRepo, rootDomainRepo, filterRepo)
	//reviewUsecase := usecases.NewReviewUsecase(domainReviewRepo, rootDomainRepo)
	//
	//baseHandler := rest.NewBaseHandler(os.Getenv("MODE"), log.Logger)
	//domainHandler := rest.NewDomainHandler(*baseHandler, domainUsecase)
	//filterHandler := rest.NewFilterHandler(*baseHandler, filterUsecase)
	//inspectHandler := rest.NewInspectHandler(*baseHandler, inspectUsecase)
	//reviewHandler := rest.NewReviewHandler(*baseHandler, reviewUsecase)
	//
	//accessTokenService := OAuthService.NewAccessTokenService(os.Getenv("ACCESS_SECRET"))
	//
	//oauthMiddleware := middleware.NewOAuthMiddleware(accessTokenService)
	//basicAuthMiddleware := middleware.NewBasicAuthMiddleware(os.Getenv("BASIC_AUTH_DOCS_USERNAME"), os.Getenv("BASIC_AUTH_DOCS_PASSWORD"))

	//checkmailServer := RPCServer.NewCheckmailServer(rpcPort, inspectUsecase)

	//app := NewConfig(*domainHandler, *filterHandler, *inspectHandler, *reviewHandler, oauthMiddleware, basicAuthMiddleware)
	app := InitializeApp()
	e := app.NewHTTPServer()
	//middleware.AddLog(e, log.Logger)
	//
	errChan := make(chan error)

	//go func() {
	//	log.Infof("starting checkmail-service RPC server on port %d\n", rpcPort)
	//	errChan <- rpc.Register(checkmailServer)
	//	errChan <- checkmailServer.Listen()
	//}()

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
