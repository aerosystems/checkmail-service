package main

import (
	"context"
	"golang.org/x/sync/errgroup"
)

// @title Checkmail Service
// @version 1.0.7
// @description A part of microservice infrastructure, who responsible for store and check email domains in black/whitelists

// @contact.name Artem Kostenko
// @contact.url https://github.com/aerosystems

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
	app := InitApp()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	group, ctx := errgroup.WithContext(ctx)

	switch app.cfg.Proto {
	case "http":
		group.Go(func() error {
			return app.httpServer.Run()
		})
	case "grpc":
		group.Go(func() error {
			return app.grpcServer.Run()
		})
	default:
		app.log.Fatalf("unknown protocol: %s", app.cfg.Proto)
	}

	group.Go(func() error {
		return app.handleSignals(ctx, cancel)
	})

	if err := group.Wait(); err != nil {
		app.log.Errorf("error occurred: %v", err)
	}
}
