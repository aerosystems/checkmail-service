package HTTPServer

import (
	"github.com/labstack/echo/v4"
)

const xAPIHeaderName = "X-Api-Key"

type ApiKeyAuth struct {
	accessUsecase AccessUsecase
}

func NewApiKeyAuth(accessUsecase AccessUsecase) *ApiKeyAuth {
	return &ApiKeyAuth{accessUsecase: accessUsecase}
}

func (a ApiKeyAuth) Auth() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			_, err := a.accessUsecase.GetAccess(c.Request().Context(), getAPIKeyFromContext(c))
			if err != nil {
				return err
			}
			return next(c)
		}
	}
}

func getAPIKeyFromContext(c echo.Context) string {
	return c.Request().Header.Get(xAPIHeaderName)
}
