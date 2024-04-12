package crypto

import (
	"strconv"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type JWTConfig struct {
	Secret    []byte
	ExpiresAt time.Duration
}

func GenerateJWT(config JWTConfig, subject int) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.RegisteredClaims{
		Subject:   strconv.FormatInt(int64(subject), 10),
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(config.ExpiresAt)),
	})

	tokenString, err := token.SignedString(config.Secret)

	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func VerifyJWT(config JWTConfig, jwtString string) (*jwt.Token, error) {
	token, err := jwt.Parse(jwtString, func(token *jwt.Token) (interface{}, error) {
		return config.Secret, nil
	}, jwt.WithValidMethods([]string{jwt.SigningMethodHS256.Name}))

	if err != nil {
		return nil, err
	}

	return token, nil
}
