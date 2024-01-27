package middleware

type TokenService interface {
	GetAccessSecret() string
	DecodeAccessToken(tokenString string) (*AccessTokenClaims, error)
}
