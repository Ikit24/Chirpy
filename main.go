package main

import (
	"log"
	"net/http"
)

type apiHandler struct{}
func (apiHandler) ServeHTTP(http.ResponseWriter, *http.Request) {}

func main() {
	mux := http.NewServeMux()
	server := http.Server{
		Addr: ":8080",
		Handler: mux,
	}
	
	mux.Handle("/", http.FileServer(http.Dir(".")))
	
	log.Fatal(server.ListenAndServe())
}
