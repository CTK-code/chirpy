package auth

import (
	"errors"
	"net/http"
	"strings"
)

func GetAPIKey(headers http.Header) (string, error) {
	key := headers.Get("Authorization")
	if key == "" {
		return "", errors.New("missing authorization header")
	}
	strippedKey, found := strings.CutPrefix(key, "ApiKey ")
	if !found {
		return "", errors.New("misformed authorization header")
	}
	return strippedKey, nil
}
