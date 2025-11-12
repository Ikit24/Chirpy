package auth

import (
	"time"
	"net/http"
	"strings"
	"errors"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

func MakeJWT(userID uuid.UUID, tokenSecret string, expiresIn time.Duration) (string, error) {
	claims := jwt.RegisteredClaims{
		Issuer:    "chirpy",
		IssuedAt:  jwt.NewNumericDate(time.Now().UTC()),
		ExpiresAt: jwt.NewNumericDate(time.Now().UTC().Add(expiresIn)),
		Subject:   userID.String(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signed, err := token.SignedString([]byte(tokenSecret))
	if err != nil {
		return "", err
	}
	return signed, nil
}

func GetBearerToken(headers http.Header) (string, error) {
	authVals := headers["Authorization"]
	if len(authVals) == 0 {
		return "", errors.New("authorization header missing")
	}

	authHeader := authVals[0]

	if !strings.HasPrefix(authHeader, "Bearer ") { return "", errors.New("invalid bearer token") }
	token := strings.TrimSpace(strings.TrimPrefix(authHeader, "Bearer "))
	if token == "" {
		return "", errors.New("invalid bearer token")
	}
	return token, nil
}
