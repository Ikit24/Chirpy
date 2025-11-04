package auth

import (
	"testing"
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
