package utils

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// MinSecretBytes is the minimum HS256 secret length (256 bits).
// Spring's NimbusJwtDecoder rejects shorter keys too.
const MinSecretBytes = 32

// TokenConfig holds everything needed to sign and verify tokens.
// It is built once at startup (see config.go) and passed in, so the
// handlers never read environment variables per request.
type TokenConfig struct {
	Secret   []byte
	Issuer   string // "iss": who created the token
	Audience string // "aud": who the token is meant for (the backend API)
	TTL      time.Duration
}

// GenerateToken creates a signed HS256 JWT for the given username.
func GenerateToken(cfg TokenConfig, username string) (string, error) {
	now := time.Now()
	claims := jwt.MapClaims{
		"sub":      username,
		"username": username,
		"iss":      cfg.Issuer,
		"aud":      cfg.Audience,
		"iat":      now.Unix(),
		"exp":      now.Add(cfg.TTL).Unix(),
	}

	// Header + "." + Payload are signed with HMAC-SHA256 using the shared secret.
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(cfg.Secret)
}

// ValidateToken verifies the signature, algorithm, expiry, issuer and audience.
// The backend does its own verification; this is used by tests and any future
// endpoint that needs to check a token.
func ValidateToken(cfg TokenConfig, tokenString string) (jwt.MapClaims, error) {
	token, err := jwt.Parse(tokenString,
		func(token *jwt.Token) (interface{}, error) {
			return cfg.Secret, nil
		},
		jwt.WithValidMethods([]string{jwt.SigningMethodHS256.Alg()}), // reject "none" / other algs
		jwt.WithIssuer(cfg.Issuer),
		jwt.WithAudience(cfg.Audience),
		jwt.WithExpirationRequired(),
	)
	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		return claims, nil
	}
	return nil, fmt.Errorf("invalid token")
}
