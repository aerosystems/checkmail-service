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
