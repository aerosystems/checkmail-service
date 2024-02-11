package HTTPServer

import OAuthService "github.com/aerosystems/checkmail-service/pkg/oauth"

type TokenService interface {
	GetAccessSecret() string
	DecodeAccessToken(tokenString string) (*OAuthService.AccessTokenClaims, error)
}
