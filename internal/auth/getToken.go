package auth

import "net/http"

func GetBearerToken(headers http.Header) (string, error) {
	TOKEN_STRING := headers.Get("Authorization")
	return
}
