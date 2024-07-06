package middleware

import (
	"github.com/labstack/echo/v4"
)

type ApiKeyAuth struct {
	accessUsecase AccessUsecase
}

func NewApiKeyAuth(accessUsecase AccessUsecase) *ApiKeyAuth {
	return &ApiKeyAuth{accessUsecase: accessUsecase}
}

func (a ApiKeyAuth) Auth() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			_, err := a.accessUsecase.GetAccess(c.Request().Header.Get("X-Api-Key"))
			if err != nil {
				return err
			}
			return next(c)
		}
	}
}
