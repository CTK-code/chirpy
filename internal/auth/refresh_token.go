package auth

import (
	"crypto/rand"
	"encoding/hex"
)

func MakeRefreshToken() (string, error) {
	key := make([]byte, 32)
	// Documentation specifies that errors are never generated
	rand.Read(key)
	token := hex.EncodeToString(key)
	return token, nil
}
