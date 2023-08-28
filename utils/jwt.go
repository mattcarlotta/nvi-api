package utils

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt"
)

var JWT_SECRET_KEY = []byte(GetEnv("JWT_SECRET_KEY"))

type JWTTokenClaim struct {
	Email string `json:"email"`
	jwt.StandardClaims
}

type JWTSessionClaim struct {
	Email  string `json:"email"`
	Name   string `json:"name"`
	UserID string `json:"userID"`
	jwt.StandardClaims
}

func GenerateUserToken(email string) (string, time.Time, error) {
	exp := time.Now().Add(time.Hour * 24 * 30)
	claims := JWTTokenClaim{
		Email: email,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: exp.Unix(),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(JWT_SECRET_KEY)
	return tokenString, exp, err
}

func ValidateUserToken(userToken string) (*JWTTokenClaim, error) {
	if len(userToken) == 0 {
		return nil, errors.New("no token was provided")
	}

	token, err := jwt.ParseWithClaims(
		userToken,
		&JWTTokenClaim{},
		func(_ *jwt.Token) (interface{}, error) {
			return JWT_SECRET_KEY, nil
		},
	)

	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(*JWTTokenClaim)
	if !ok {
		return nil, errors.New("unable to parse user token")
	} else if claims.ExpiresAt < time.Now().Local().Unix() {
		return nil, errors.New("token expired")
	}

	return claims, nil
}

func ValidateSessionToken(jwtCookie string) (*JWTSessionClaim, error) {
	if len(jwtCookie) == 0 {
		return nil, errors.New("you must be logged in order to do that")
	}

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
		return nil, errors.New("unable to parse session token")
	} else if claims.ExpiresAt < time.Now().Local().Unix() {
		return nil, errors.New("session expired")
	}

	return claims, nil
}
