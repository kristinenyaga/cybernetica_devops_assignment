package main

import (
	"log"
	"net/http"
	"os"
	"time"
	"context"
	"os/signal"
	"syscall"
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
	log.SetFlags(0)

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
	
	server := &http.Server{
		Addr:    ":" + port,
		Handler: loggedMux,
	}

	// Starting server in a goroutine
	go func() {
		log.Printf("Server is starting on port %s", port)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("ListenAndServe(): %v", err)
		}
	}()

	// Creating a channel to listen for OS signals
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit // blocks here until signal arrives

	log.Println("Shutdown requested, shutting down server...")

	// Context with timeout for graceful shutdown
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}

	log.Println("Server gracefully shut down")
	
}