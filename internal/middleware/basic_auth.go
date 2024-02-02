package middleware

import (
	"github.com/labstack/echo/v4"
	"net/http"
	"os"
)

type BasicAuthMiddleware interface {
	BasicAuthMiddleware(next echo.HandlerFunc) echo.HandlerFunc
}

type BasicAuthMiddlewareImpl struct {
	username string
	password string
}

func NewBasicAuthMiddleware() *BasicAuthMiddlewareImpl {
	return &BasicAuthMiddlewareImpl{
		username: os.Getenv("BASIC_AUTH_DOCS_USERNAME"),
		password: os.Getenv("BASIC_AUTH_DOCS_PASSWORD"),
	}
}

func (b *BasicAuthMiddlewareImpl) BasicAuthMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		username, password, ok := c.Request().BasicAuth()

		if !ok || !b.checkCredentials(username, password) {
			c.Response().Header().Set("WWW-Authenticate", `Basic realm="Restricted"`)
			return echo.NewHTTPError(http.StatusUnauthorized, "Invalid credentials")
		}

		return next(c)
	}
}

func (b *BasicAuthMiddlewareImpl) checkCredentials(username, password string) bool {
	return username == b.username && password == b.password
}
