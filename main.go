package main

import (
	"net/http"
	"os"
	"time"
	"log"
)
var (
	allowOrigin     string
	responseMessage string
)

func homeHandler(w http.ResponseWriter, _ *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", allowOrigin)
	w.Write([]byte(responseMessage))
}

func healthHandler(w http.ResponseWriter, _ *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("ok"))
}

func readyHandler(w http.ResponseWriter, _ *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("ready"))
}

type statusRecorder struct {
	http.ResponseWriter
	statusCode int
}

func (r *statusRecorder) WriteHeader(code int) {
	r.statusCode = code
	r.ResponseWriter.WriteHeader(code)
}

func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()                      
		rec := &statusRecorder{
			ResponseWriter: w,
			statusCode: http.StatusOK,
		}

		next.ServeHTTP(rec, r)                      
		duration := time.Since(start)             
		log.Printf(
		`{"time":"%s","method":"%s","path":"%s","status":%d,"duration":"%s"}`,
		time.Now().Format(time.RFC3339),
		r.Method,
		r.URL.Path,
		rec.statusCode,
		duration,
		) 
	})
}

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8000"
	}
	responseMessage = os.Getenv("RESPONSE_MESSAGE")
	if responseMessage == "" {
		responseMessage = "Service request succeeded!"
	}
	allowOrigin = os.Getenv("ALLOW_ORIGIN")
	if allowOrigin == "" {
		allowOrigin = "*"
	}
	
	mux := http.NewServeMux()

	mux.HandleFunc("/", homeHandler)
	mux.HandleFunc("/health", healthHandler)
	mux.HandleFunc("/ready", readyHandler)

	loggedMux := loggingMiddleware(mux)
	
	log.Fatal(http.ListenAndServe(":"+port, loggedMux))
}