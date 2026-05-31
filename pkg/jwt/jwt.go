package jwt

import (
	"time"

	"github.com/golang-jwt/jwt"
)

func CreateToken(id int64, username, secretKey string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"id":       id,
		"username": username,
		"exp":      time.Now().Add(60 * time.Minute).Unix(),
	})

	key := []byte(secretKey)
	tokenStr, err := token.SignedString(key)

	return tokenStr, err
}
