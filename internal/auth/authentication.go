package main

import (
	"log"

	"github.com/alexedwards/argon2id"
)

func HashPassword(password strnig) (string, error) {
	hash, err := argon2id.CreateHAsh(password, argon2id.DefaultParams)
	if err != nil {
		log.Fatal(err)
	}
}
