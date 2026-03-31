package main

import (
	"net/http"
	"os"
)

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8000"
	}

	responseMessage := os.Getenv("RESPONSE_MESSAGE")
	if responseMessage == "" {
		responseMessage = "Service request succeeded!"
	}

	allowOrigin := os.Getenv("ALLOW_ORIGIN")
	if allowOrigin == "" {
		allowOrigin = "*"
	}

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", allowOrigin)
		w.Write([]byte(responseMessage))
	})

	http.ListenAndServe(":"+port, nil)
}