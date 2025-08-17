package utils

import (
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type Payload struct {
	Id  string `json:"id"`
	nbf string
	exp string
	jwt.RegisteredClaims
}

const Secret string = "JWTSecret123"

func CreateToken(id string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"id":  id,
		"nbf": time.Date(2015, 10, 10, 12, 0, 0, 0, time.UTC).Unix(),
		"exp": time.Now().Add(time.Hour * 24).Unix(),
	})

	tokenString, err := token.SignedString(StringToBytes(Secret))
	fmt.Println(err)
	if err != nil {
		return "", errors.New("could not create jwt")
	}

	return tokenString, nil
}

func ValidateToken(tokenString string) (*Payload, error) {
	token, err := jwt.ParseWithClaims(tokenString, &Payload{}, func(token *jwt.Token) (any, error) {
		return StringToBytes(Secret), nil
	})

	if err != nil {
		return nil, errors.New("not a valid jwt")
	}

	claims := token.Claims.(*Payload)

	return claims, nil
}
