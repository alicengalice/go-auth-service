package utils

import (
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

func cfg() TokenConfig {
	return TokenConfig{
		Secret:   []byte("0123456789abcdef0123456789abcdef"),
		Issuer:   "test-issuer",
		Audience: "test-audience",
		TTL:      time.Hour,
	}
}

func TestGenerateToken_HasExpectedClaims(t *testing.T) {
	tok, err := GenerateToken(cfg(), "admin")
	if err != nil {
		t.Fatal(err)
	}
	claims, err := ValidateToken(cfg(), tok)
	if err != nil {
		t.Fatalf("validate: %v", err)
	}
	for _, k := range []string{"sub", "username", "iss", "aud", "iat", "exp"} {
		if _, ok := claims[k]; !ok {
			t.Errorf("missing claim %q", k)
		}
	}
}

func TestValidateToken_Rejects(t *testing.T) {
	good := cfg()

	wrongSecret := good
	wrongSecret.Secret = []byte("ffffffffffffffffffffffffffffffff")

	wrongIssuer := good
	wrongIssuer.Issuer = "someone-else"

	wrongAudience := good
	wrongAudience.Audience = "other-api"

	expired := good
	expired.TTL = -time.Minute

	cases := []struct {
		name     string
		signWith TokenConfig
	}{
		{"wrong secret", wrongSecret},
		{"wrong issuer", wrongIssuer},
		{"wrong audience", wrongAudience},
		{"expired", expired},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			tok, err := GenerateToken(c.signWith, "admin")
			if err != nil {
				t.Fatal(err)
			}
			if _, err := ValidateToken(good, tok); err == nil {
				t.Error("expected validation to fail")
			}
		})
	}
}

func TestValidateToken_RejectsNoneAlgorithm(t *testing.T) {
	tok := jwt.NewWithClaims(jwt.SigningMethodNone, jwt.MapClaims{
		"sub": "admin", "iss": "test-issuer", "aud": "test-audience",
		"exp": time.Now().Add(time.Hour).Unix(),
	})
	s, err := tok.SignedString(jwt.UnsafeAllowNoneSignatureType)
	if err != nil {
		t.Fatal(err)
	}
	if _, err := ValidateToken(cfg(), s); err == nil {
		t.Error("alg=none token must be rejected")
	}
}
