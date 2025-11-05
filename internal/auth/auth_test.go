package auth

import (
	"testing"
	"time"

	"github.com/google/uuid"
)

func TestCheckPasswordHashing(t *testing.T) {
	password1 := "correctPassword123!"
	password2 := "anotherPassword456!"
	hash1, _ := HashPassword(password1)
	hash2, _ := HashPassword(password2)

	tests := []struct {
		name          string
		password      string
		hash          string
		wantErr       bool
		matchPassword bool
	}{
		{
			name:		  "Correct password",
			password:	  password1,
			hash:		  hash1,
			wantErr:	  false,
			matchPassword: true,
		},
		{
			name:		  "Incorrect password",
			password:	  "wrongPassword",
			hash:		  hash1,
			wantErr:	  false,
			matchPassword: false,
		},
		{
			name:		  "Password doesn't match different hash",
			password:	  password1,
			hash:		  hash2,
			wantErr:	  false,
			matchPassword: false,
		},
		{
			name:		  "Empty password",
			password:	  "",
			hash:		  hash1,
			wantErr:	  false,
			matchPassword: false,
		},
		{
			name:		  "Invalid hash",
			password:	  password1,
			hash:		  "invalidhash",
			wantErr:	  true,
			matchPassword: false,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			match, err := CheckPasswordHash(tc.password, tc.hash)
			if (err != nil) != tc.wantErr {
				t.Errorf("CheckPasswordHash() error = %v, wantErr %v", err, tc.wantErr)
			}
			if !tc.wantErr && match != tc.matchPassword {
				t.Errorf("CheckPasswordHash() expects %v, got %v", tc.matchPassword, match)
			}
		})
	}
}

func TestMakeAndValidateJWT_Valid(t *testing.T) {
	tests := []struct {
		name		string
		expiresIn	time.Duration
		makeWith	string
		checkWith	string
		wantErr		bool
	}{
		{
			name:		"valid",
			expiresIn:	time.Minute,
			makeWith:	"secretA",
			checkWith:	"secretA",
			wantErr:	false,
		},
		{
			name:		"expired",
			expiresIn:	-time.Minute,
			makeWith:	"secretA",
			checkWith:	"secretA",
			wantErr:	true,
		},
		{
			name:		"wrong-secret",
			expiresIn:	time.Minute,
			makeWith:	"secretA",
			checkWith:	"secretB",
			wantErr:	true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			userID := uuid.New()
			token, err := MakeJWT(userID, tc.makeWith, tc.expiresIn)
			if err != nil {
				t.Errorf("Could create JWT: %v", err)
			}

			valID, err := ValidateJWT(token, tc.checkWith)
			if tc.wantErr {
				if err == nil {
					t.Fatalf("expected error, got none")
				}
				return
			}
			if err != nil {
				t.Errorf("unexpected error: %v", err)
			}
			if valID != userID {
				t.Errorf("want %s, got %s", userID, valID)
			}
		})
	}
}
