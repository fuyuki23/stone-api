package token

import (
	"crypto/rsa"
	"stone-api/internal/config"
	"stone-api/internal/model"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/rs/zerolog/log"
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

func ValidateToken(tokenType string, userToken string) (bool, error) {
	claims, err := getClaims(userToken)
	if err != nil {
		return false, err
	}
	if claims == nil {
		return false, nil
	}

	if claims.Issuer != ISSUER {
		log.Warn().Msg("invalid issuer")
		return false, nil
	}
	if claims.Type != tokenType {
		log.Warn().Msg("invalid token type")
		return false, nil
	}
	if claims.Email == "" {
		log.Warn().Msg("invalid email")
		return false, nil
	}

	return true, nil
}

func GetEmailFromToken(token string) (string, error) {
	claims, err := getClaims(token)
	if err != nil {
		return "", err
	}
	if claims == nil {
		return "", model.ErrUnauthorized
	}

	return claims.Email, nil
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

func getClaims(userToken string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(userToken, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
			log.Debug().Msg("invalid signing method")
			return nil, jwt.ErrInvalidKeyType
		}

		return &(getKey().PublicKey), nil
	})
	if err != nil {
		return nil, err
	}
	if !token.Valid {
		return nil, jwt.ErrTokenExpired
	}

	if claims, ok := token.Claims.(*Claims); ok {
		return claims, nil
	}

	log.Warn().Msg("invalid claims")
	return nil, model.ErrUnknown
}
