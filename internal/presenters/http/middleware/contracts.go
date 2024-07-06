package middleware

import "github.com/aerosystems/checkmail-service/internal/models"

type AccessUsecase interface {
	GetAccess(apiKey string) (*models.Access, error)
}
