package auth

import (
	"crypto/rand"
	"encoding/hex"
)

func MakeRefreshToken() (string, error) {
	dat := make([]byte, 32)
	rand.Read(dat)
	code := hex.EncodeToString(dat)
	return code, nil
}
