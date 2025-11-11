package auth

import (
	"errors"
	"fmt"
	"net/http"
	"strings"
)

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
