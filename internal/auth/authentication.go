package auth

import (
	"errors"
	"strings"
	"net/http"

	"github.com/alexedwards/argon2id"
)

func HashPassword(password string) (string, error) {
	hash, err := argon2id.CreateHash(password, argon2id.DefaultParams)
	if err != nil {
		return "", err
	}
	return hash, nil
}

func CheckPasswordHash(password, hash string) (bool, error) {
	match, err := argon2id.ComparePasswordAndHash(password, hash)
	if err != nil {
		return false, err
	}
	return match, nil
}

func GetAPIKey(headers http.Header) (string, error) {
	authHeader := headers.Get("Authorization")
	if len(authHeader) == 0 {
		return "", errors.New("authorization header missing")
	}

	if !strings.HasPrefix(authHeader, "ApiKey ") {
		return "", errors.New("invalid API key")
	}

	apiKey := strings.TrimSpace(strings.TrimPrefix(authHeader, "ApiKey "))
	if apiKey == "" {
		return "", errors.New("invalid API key")
	}
	return apiKey, nil
}
