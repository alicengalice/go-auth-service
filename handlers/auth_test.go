package handlers

import (
	"encoding/json"
	"go-auth-service/models"
	"go-auth-service/utils"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
)

func testConfig() AuthConfig {
	return AuthConfig{
		AdminUsername: "admin",
		AdminPassword: "s3cret-pass",
		Token: utils.TokenConfig{
			Secret:   []byte("0123456789abcdef0123456789abcdef"),
			Issuer:   "test-issuer",
			Audience: "test-audience",
			TTL:      time.Hour,
		},
	}
}

func doLogin(t *testing.T, method, body string) (*httptest.ResponseRecorder, models.LoginResponse) {
	t.Helper()
	req := httptest.NewRequest(method, "/login", strings.NewReader(body))
	rec := httptest.NewRecorder()
	Login(testConfig())(rec, req)

	var resp models.LoginResponse
	if rec.Header().Get("Content-Type") == "application/json" {
		if err := json.NewDecoder(rec.Body).Decode(&resp); err != nil {
			t.Fatalf("decode response: %v", err)
		}
	}
	return rec, resp
}

func TestLogin_ValidCredentials_ReturnsVerifiableToken(t *testing.T) {
	rec, resp := doLogin(t, http.MethodPost, `{"username":"admin","password":"s3cret-pass"}`)

	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d, want 200", rec.Code)
	}
	if resp.Token == "" {
		t.Fatal("expected a token")
	}

	claims, err := utils.ValidateToken(testConfig().Token, resp.Token)
	if err != nil {
		t.Fatalf("token should validate: %v", err)
	}
	if claims["sub"] != "admin" || claims["iss"] != "test-issuer" {
		t.Errorf("unexpected claims: %v", claims)
	}
}

func TestLogin_WrongCredentials_Returns401(t *testing.T) {
	cases := map[string]string{
		"wrong password": `{"username":"admin","password":"nope"}`,
		"wrong username": `{"username":"root","password":"s3cret-pass"}`,
		"empty":          `{"username":"","password":""}`,
		"prefix of pass": `{"username":"admin","password":"s3cret"}`,
	}
	for name, body := range cases {
		t.Run(name, func(t *testing.T) {
			rec, resp := doLogin(t, http.MethodPost, body)
			if rec.Code != http.StatusUnauthorized {
				t.Fatalf("status = %d, want 401", rec.Code)
			}
			if resp.Token != "" {
				t.Error("no token should be returned")
			}
		})
	}
}

func TestLogin_InvalidJSON_Returns400(t *testing.T) {
	rec, _ := doLogin(t, http.MethodPost, `not json`)
	if rec.Code != http.StatusBadRequest {
		t.Fatalf("status = %d, want 400", rec.Code)
	}
}

func TestLogin_NonPost_Returns405(t *testing.T) {
	rec, _ := doLogin(t, http.MethodGet, "")
	if rec.Code != http.StatusMethodNotAllowed {
		t.Fatalf("status = %d, want 405", rec.Code)
	}
}
