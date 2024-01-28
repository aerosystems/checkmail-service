package middleware

import OAuthService "github.com/aerosystems/checkmail-service/pkg/oauth_service"

type TokenService interface {
	GetAccessSecret() string
	DecodeAccessToken(tokenString string) (*OAuthService.AccessTokenClaims, error)
}
