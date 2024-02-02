package OAuthService

import (
	"github.com/golang-jwt/jwt"
	"os"
)

type AccessTokenClaims struct {
	AccessUuid string `json:"accessUuid"`
	UserUuid   string `json:"userUuid"`
	UserRole   string `json:"userRole"`
	Exp        int    `json:"exp"`
	jwt.StandardClaims
}

type AccessTokenService struct {
	accessSecret string
}

func NewAccessTokenService() *AccessTokenService {
	return &AccessTokenService{
		accessSecret: os.Getenv("ACCESS_SECRET"),
	}
}

func (r AccessTokenService) GetAccessSecret() string {
	return r.accessSecret
}

func (r AccessTokenService) DecodeAccessToken(tokenString string) (*AccessTokenClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &AccessTokenClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(r.accessSecret), nil
	})
	if claims, ok := token.Claims.(*AccessTokenClaims); ok && token.Valid {
		return claims, nil
	} else {
		return nil, err
	}
}
