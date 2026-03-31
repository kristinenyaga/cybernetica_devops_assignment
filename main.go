package main

import (
	"net/http"
	"os"
	"time"
	"log"
)

func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()                      
		next.ServeHTTP(w, r)                      
		duration := time.Since(start)             
		log.Printf("%s %s %v", r.Method, r.URL.Path, duration) 
	})
}

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

	http.Handle("/", loggingMiddleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", allowOrigin)
		w.Write([]byte(responseMessage))
	})))

	http.Handle("/health", loggingMiddleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("ok"))
	})))

	http.Handle("/ready", loggingMiddleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("ready"))
	})))


	http.ListenAndServe(":"+port, nil)
}