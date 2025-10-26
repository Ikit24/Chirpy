package main

import (
	"encoding/json"
	"net/http"
	"log"
)

type chirps struct {
	Body string `json:"body"`
		Id string
	}

	func (cfg *apiConfig) handlerChirps(w http.ResponseWrite, r *http.Rewquest) {

	}
