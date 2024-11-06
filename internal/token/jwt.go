package token

import (
	"crypto/rsa"
	"github.com/golang-jwt/jwt/v5"
	"stone-api/internal/config"
	"stone-api/internal/model"
	"time"
)

type Claims struct {
	jwt.RegisteredClaims
	Type  string `json:"type"`
	Email string `json:"email"`
}

const ISSUER = "stone-api"

var privateKey *rsa.PrivateKey

func CreateTokens(user model.User) (*model.Tokens, error) {
	accessToken, err := createAccessToken(&user)
	if err != nil {
		return nil, err
	}
	refreshToken, err := createRefreshToken(&user)
	if err != nil {
		return nil, err
	}

	return &model.Tokens{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}

func createAccessToken(user *model.User) (string, error) {
	return jwt.NewWithClaims(jwt.SigningMethodRS256, Claims{
		Type:  "access",
		Email: user.Email,
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    ISSUER,
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Minute * 15)),
		},
	}).SignedString(getKey())
}

func createRefreshToken(user *model.User) (string, error) {
	return jwt.NewWithClaims(jwt.SigningMethodRS256, Claims{
		Type:  "access",
		Email: user.Email,
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    ISSUER,
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Minute * 15)),
		},
	}).SignedString(getKey())
}

func getKey() *rsa.PrivateKey {
	if privateKey == nil {
		key, err := jwt.ParseRSAPrivateKeyFromPEM([]byte(config.Get().Server.Jwt.PrivateKey))
		if err != nil {
			panic(err)
		}
		privateKey = key
	}

	return privateKey
}
