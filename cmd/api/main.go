package main

import (
	"fmt"
	"github.com/aerosystems/checkmail-service/internal/handlers"
	"github.com/aerosystems/checkmail-service/internal/models"
	"github.com/aerosystems/checkmail-service/internal/repository"
	"github.com/aerosystems/checkmail-service/pkg/mygorm"

	"log"
	"net/http"
)

const webPort = "80"

type Config struct {
	BaseHandler *handlers.BaseHandler
}

// @title Checkmail Service
// @version 1.0
// @description A part of microservice infrastructure, who responsible for store and check email domains in black/whitelists

// @contact.name Artem Kostenko
// @contact.url https://github.com/aerosystems

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html

// @host localhost:8083
// @BasePath /v1
func main() {
	clientGORM := mygorm.NewClient()
	if err := clientGORM.AutoMigrate(&models.Domain{}, &models.RootDomain{}); err != nil {
		log.Panic(err)
	}
	domainRepo := repository.NewDomainRepo(clientGORM)
	rootDomainRepo := repository.NewRootDomainRepo(clientGORM)

	app := Config{
		BaseHandler: handlers.NewBaseHandler(domainRepo, rootDomainRepo),
	}

	srv := &http.Server{
		Addr:    fmt.Sprintf(":%s", webPort),
		Handler: app.routes(),
	}

	log.Printf("Starting authentication end service on port %s\n", webPort)
	err := srv.ListenAndServe()

	if err != nil {
		log.Panic(err)
	}
}
