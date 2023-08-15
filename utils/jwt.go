package utils

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt"
)

var JWT_SECRET_KEY = []byte(GetEnv("JWT_SECRET_KEY"))

type JWTSessionClaim struct {
	Email  string `json:"email"`
	Name   string `json:"name"`
	UserId string `json:"userId"`
	jwt.StandardClaims
}

func ValidateSessionToken(jwtCookie string) (*JWTSessionClaim, error) {
	token, err := jwt.ParseWithClaims(
		jwtCookie,
		&JWTSessionClaim{},
		func(_ *jwt.Token) (interface{}, error) {
			return JWT_SECRET_KEY, nil
		},
	)

	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(*JWTSessionClaim)
	if !ok {
		return nil, errors.New("Unable to parse session token.")
	} else if claims.ExpiresAt < time.Now().Local().Unix() {
		return nil, errors.New("Session expired.")
	}

	return claims, nil
}
