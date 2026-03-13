package main

import (
	"fmt"
	"go-auth-service/handlers"
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
)

// corsMiddleware adds CORS headers to allow browser requests from the frontend
func corsMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "http://localhost:5173")
		w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

		// Handle preflight OPTIONS request
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}

		next(w, r)
	}
}

func main() {
	// Load .env file
	err := godotenv.Load()
	if err != nil {
		log.Println("Warning: .env file not found, using system environment variables")
	}

	// Get port from environment (default to 8081)
	port := os.Getenv("PORT")
	if port == "" {
		port = "8081"
	}

	// Setup routes
	http.HandleFunc("/login", corsMiddleware(handlers.Login))
	http.HandleFunc("/health", corsMiddleware(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, "Auth service is running")
	}))

	// Start server
	fmt.Printf("Starting auth service on : %s\n", port)
	err = http.ListenAndServe(":"+port, nil)
	if err != nil {
		panic(err)
	}
}
