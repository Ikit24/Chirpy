package main

import (
	//"fmt"
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

// 	mux.Handle("/api/", apiHandler{})
	//mux.HandleFunc("/", func(w http.ResponseWriter, req *http.Request) {
	//	if req.URL.Path != "/" {
	//		http.NotFound(w, req)
	//		return
	//	}
		//fmt.Fprintf(w, "Welcome to the homepage!")
	//})
	
	server.ListenAndServe()
}
