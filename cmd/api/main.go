package main

import (
	"fmt"
	"github.com/aerosystems/checkmail-service/internal/handlers"
	"github.com/aerosystems/checkmail-service/internal/models"
	"github.com/aerosystems/checkmail-service/internal/repository"
	"github.com/aerosystems/checkmail-service/pkg/gorm_postgres"

	"log"
	"net/http"
)

// @title Checkmail Service
// @version 1.0.6
// @description A part of microservice infrastructure, who responsible for store and check email domains in black/whitelists

// @contact.name Artem Kostenko
// @contact.url https://github.com/aerosystems

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html

// @securityDefinitions.apikey X-API-KEY
// @in header
// @name X-API-KEY
// @description Should contain Token, digits and letters, 64 symbols length

// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @description Should contain Access JWT Token, with the Bearer started

// @host localhost:8083
// @BasePath /
func main() {
	clientGORM := mygorm.NewClient()
	if err := clientGORM.AutoMigrate(&models.Domain{}, &models.RootDomain{}); err != nil {
		log.Panic(err)
	}
	domainRepo := repository.NewDomainRepo(clientGORM)
	rootDomainRepo := repository.NewRootDomainRepo(clientGORM)

	app := Config{
		WebPort:     "80",
		BaseHandler: handlers.NewBaseHandler(domainRepo, rootDomainRepo),
	}

	srv := &http.Server{
		Addr:    fmt.Sprintf(":%s", app.WebPort),
		Handler: app.routes(),
	}

	log.Printf("Starting authentication end service on port %s\n", app.WebPort)
	err := srv.ListenAndServe()

	if err != nil {
		log.Panic(err)
	}
}
