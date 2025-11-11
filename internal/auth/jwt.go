package auth

import (
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

func MakeJWT(userID uuid.UUID, tokenSecret string, expiresIn time.Duration) (string, error) {
	claims := jwt.RegisteredClaims{
		Issuer:    "chirpy",
		ExpiresAt: jwt.NewNumericDate(time.Now().UTC().Add(expiresIn)),
		Subject:   userID.String(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(tokenSecret))
	if err != nil {
		return "", errors.New("couldn't sign JWT token: " + err.Error())
	}
	return tokenString, nil
}

func ValidateJWT(tokenString, tokenSecret string) (uuid.UUID, error) {
	type myClaims struct {
		Foo string `json:"foo"`
		jwt.RegisteredClaims
	}

	token, err := jwt.ParseWithClaims(
		tokenString,
		&myClaims{},
		func(token *jwt.Token) (any, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, errors.New("unexpected signing method")
			}
			return []byte(tokenSecret), nil
		})
	if err != nil {
		return uuid.Nil, errors.New("couldn't parse JWT token: " + err.Error())
	}
	if token == nil || !token.Valid {
		return uuid.Nil, errors.New("invalid JWT token")
	}

	claims, ok := token.Claims.(*myClaims)
	if !ok {
		return uuid.Nil, errors.New("invalid JWT token")
	}

	userID, err := uuid.Parse(claims.Subject)
	if err != nil {
		return uuid.Nil, errors.New("couldn't parse user ID from JWT token: " + err.Error())
	}

	return userID, nil
}

func GetBearerToken(headers http.Header) (string, error) {
	bearerToken := headers.Get("Authorization")
	if bearerToken == "" {
		return "", errors.New("no Authorization header found")
	}
	splitToken := strings.SplitN(bearerToken, " ", 2)
	if len(splitToken) != 2 || strings.ToLower(splitToken[0]) != "bearer" {
		fmt.Println("Invalid Authorization header format")
	}
	return splitToken[1], nil
}
