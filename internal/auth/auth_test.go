package auth_test

import (
	"testing"
)

func TestHashing(t *testing.T) {
	password1 := "correctPassword123!"
	password2 := "anotherPassword456!"
	hash1, _ := HashPassword(password1)
	hash2, _ := HashPassword(password2)

	test := []struct {
		name          string
		password      string
		hash          string
		wantErr       bool
		matchPassword bool
	}
}
