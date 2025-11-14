package auth

import (
	"errors"
	"fmt"
	"net/http"
	"strings"
)

func GetAPIKey(headers http.Header) (string, error) {
	polkaToken := headers.Get("Authorization")

	fmt.Println("original token:", polkaToken)

	if polkaToken == "" {
		return "", errors.New("no Authorization header found")
	}
	splitToken := strings.SplitN(polkaToken, " ", 2)

	fmt.Println("split strings:", splitToken)

	if len(splitToken) != 2 || strings.ToLower(splitToken[0]) != "apikey" {
		return "", errors.New("polka token in bad authorization header format")
	}
	return splitToken[1], nil
}
