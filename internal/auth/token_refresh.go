package auth

import(
	"crypto/rand"
	"encoding/hex"
	"fmt"
)

func MakeRefreshToken() (string, error) {
	key := make([] byte, 32)
	check, err := rand.Read(key)
	if err != nil {
		return "", err
	}
	if check != len(key) {
		return "", fmt.Errorf("couldn't generate secure token: only got %d of 32 random bytes", check)
	}
	return hex.EncodeToString(key), nil
}
