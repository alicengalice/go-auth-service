package main

import (
	"fmt"
	"go-auth-service/handlers"
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
)

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
	http.HandleFunc("/login", handlers.Login)
	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, "Auth service is running")
	})

	// Start server
	fmt.Printf("Starting auth service on : %s\n", port)
	err = http.ListenAndServe(":"+port, nil)
	if err != nil {
		panic(err)
	}
}
