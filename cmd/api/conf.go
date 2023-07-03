package main

import (
	"github.com/aerosystems/checkmail-service/internal/handlers"
)

type Config struct {
	WebPort     string
	BaseHandler *handlers.BaseHandler
}
