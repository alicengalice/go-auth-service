package handlers

import (
	"crypto/sha256"
	"crypto/subtle"
	"encoding/json"
	"go-auth-service/models"
	"go-auth-service/utils"
	"net/http"
)

// AuthConfig holds the admin credentials and token settings, validated at startup.
type AuthConfig struct {
	AdminUsername string
	AdminPassword string
	Token         utils.TokenConfig
}

// constantTimeEqual compares two strings without leaking, through timing,
// how many leading characters matched. Hashing first makes both inputs
// the same length, because subtle.ConstantTimeCompare returns early on a length mismatch.
func constantTimeEqual(a, b string) bool {
	ha := sha256.Sum256([]byte(a))
	hb := sha256.Sum256([]byte(b))
	return subtle.ConstantTimeCompare(ha[:], hb[:]) == 1
}

// Login returns the handler for POST /login.
func Login(cfg AuthConfig) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Only accept POST
		if r.Method != http.MethodPost {
			http.Error(w, "Method not allowed. POST method only", http.StatusMethodNotAllowed)
			return
		}

		// Parse JSON request body
		var req models.LoginRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "Invalid JSON", http.StatusBadRequest)
			return
		}

		// Check both fields every time (no short-circuit), so the response
		// time doesn't reveal whether the username was right.
		userOK := constantTimeEqual(req.Username, cfg.AdminUsername)
		passOK := constantTimeEqual(req.Password, cfg.AdminPassword)
		if !(userOK && passOK) {
			writeJSON(w, http.StatusUnauthorized, models.LoginResponse{Error: "Invalid credentials"})
			return
		}

		token, err := utils.GenerateToken(cfg.Token, cfg.AdminUsername)
		if err != nil {
			http.Error(w, "Failed to generate token", http.StatusInternalServerError)
			return
		}

		writeJSON(w, http.StatusOK, models.LoginResponse{Token: token})
	}
}

func writeJSON(w http.ResponseWriter, status int, body any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(body)
}
