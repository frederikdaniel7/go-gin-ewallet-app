package utils

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

func CreateAndSign(userId int, secretKey string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"id":  userId,
		"iss": "the-issuer",
		"exp": time.Now().Add(60 * time.Minute).Unix(),
		"iat": time.Now().Unix(),
	})

	signed, err := token.SignedString([]byte(secretKey))
	if err != nil {
		return "", err
	}

	return signed, nil
}

func ParseAndVerify(signed string, secretKey string) (jwt.MapClaims, error) {
	token, err := jwt.Parse(signed, func(token *jwt.Token) (interface{}, error) {
		return []byte(secretKey), nil
	}, jwt.WithIssuer("the-issuer"),
		jwt.WithValidMethods([]string{jwt.SigningMethodHS256.Name}),
		jwt.WithExpirationRequired(),
	)
	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok {
		return claims, nil
	} else {
		return nil, fmt.Errorf("unknown claims")
	}
}
