package main

import (
	"strings"
	"testing"
)

func setValidEnv(t *testing.T) {
	t.Setenv("ADMIN_USERNAME", "admin")
	t.Setenv("ADMIN_PASSWORD", "pass")
	t.Setenv("JWT_SECRET", "0123456789abcdef0123456789abcdef")
	t.Setenv("PORT", "")
	t.Setenv("JWT_ISSUER", "")
	t.Setenv("JWT_AUDIENCE", "")
}

func TestLoadConfig_ValidEnv_UsesDefaults(t *testing.T) {
	setValidEnv(t)
	cfg, err := LoadConfig()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if cfg.Port != defaultPort || cfg.Auth.Token.Issuer != defaultIssuer || cfg.Auth.Token.Audience != defaultAudience {
		t.Errorf("defaults not applied: %+v", cfg)
	}
}

func TestLoadConfig_MissingOrWeak_ReturnsError(t *testing.T) {
	cases := map[string]struct{ key, value, wantMsg string }{
		"no username":  {"ADMIN_USERNAME", "", "ADMIN_USERNAME"},
		"no password":  {"ADMIN_PASSWORD", "", "ADMIN_PASSWORD"},
		"no secret":    {"JWT_SECRET", "", "JWT_SECRET is not set"},
		"short secret": {"JWT_SECRET", "too-short", "at least 32 bytes"},
	}
	for name, c := range cases {
		t.Run(name, func(t *testing.T) {
			setValidEnv(t)
			t.Setenv(c.key, c.value)
			_, err := LoadConfig()
			if err == nil || !strings.Contains(err.Error(), c.wantMsg) {
				t.Fatalf("err = %v, want message containing %q", err, c.wantMsg)
			}
		})
	}
}
