package main

import (
	"fmt"
	"github.com/aerosystems/checkmail-service/internal/handlers"
	"github.com/aerosystems/checkmail-service/internal/models"
	"github.com/aerosystems/checkmail-service/internal/repository"
	RPCServer "github.com/aerosystems/checkmail-service/internal/rpc_server"
	"github.com/aerosystems/checkmail-service/internal/services"
	GormPostgres "github.com/aerosystems/checkmail-service/pkg/gorm_postgres"
	"github.com/aerosystems/checkmail-service/pkg/logger"
	"github.com/sirupsen/logrus"
	"net/http"
	"net/rpc"
	"os"
)

const webPort = 80
const rpcPort = 5001

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

// @host gw.verifire.com/checkmail
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

	inspectService := services.NewInspectService(log.Logger, domainRepo, rootDomainRepo)

	checkmailServer := RPCServer.NewCheckmailServer(rpcPort, inspectService)

	baseHandler := handlers.NewBaseHandler(log.Logger, domainRepo, rootDomainRepo, filterRepo, domainReviewRepo, inspectService)

	app := Config{baseHandler}

	srv := &http.Server{
		Addr:    fmt.Sprintf(":%d", webPort),
		Handler: app.routes(log.Logger),
	}

	errChan := make(chan error)

	go func() {
		log.Infof("starting checkmail-service RPC server on port %d\n", rpcPort)
		errChan <- rpc.Register(checkmailServer)
		errChan <- checkmailServer.Listen()
	}()

	go func() {
		log.Infof("starting checkmail-service HTTP server on port %d\n", webPort)
		errChan <- srv.ListenAndServe()
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
