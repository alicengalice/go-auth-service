package handlers

import (
	"encoding/json"
	"go-auth-service/models"
	"go-auth-service/utils"
	"net/http"
	"os"
)

// GetAdminCredentials retrievees admin credentials from environment
func GetAdminCredentials() (string, string) {
	username := os.Getenv("ADMIN_USERNAME")
	password := os.Getenv("ADMIN_PASSWORD")

	if username == "" || password == "" {
		panic("ADMIN_USERNAME or ADMIN_PASSWORD environment variable is not set")
	}

	return username, password
}

// Login handles POST /login
func Login(w http.ResponseWriter, r *http.Request) {
	// Only accept POST
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed. POST method only", http.StatusMethodNotAllowed)
		return
	}

	// Parse JSON request body
	var req models.LoginRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	// Get admin credentials from environment
	adminUsername, adminPassword := GetAdminCredentials()

	// Validate credentials
	if req.Username != adminUsername || req.Password != adminPassword {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(models.LoginResponse{
			Error: "Invalid credentials",
		})
		return
	}

	// Generate JWT token
	token, err := utils.GenerateToken(req.Username)
	if err != nil {
		http.Error(w, "Failed to generate token", http.StatusInternalServerError)
		return
	}

	// Return token
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(models.LoginResponse{
		Token: token,
	})
}
