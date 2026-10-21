package config

import (
	"strings"
	"testing"
)

// Literal copy on purpose: the test must fail if the implementation constant drifts.
const knownLegacySecret = "mcpark-secret-change-in-production"

func envFrom(m map[string]string) func(string) string {
	return func(k string) string { return m[k] }
}

func TestLoadRejectsWeakSecrets(t *testing.T) {
	cases := map[string]string{
		"missing":        "",
		"legacy default": knownLegacySecret,
		"too short":      strings.Repeat("a", 31),
	}
	for name, secret := range cases {
		t.Run(name, func(t *testing.T) {
			cfg, err := LoadWithGetenv(envFrom(map[string]string{"JWT_SECRET": secret}))
			if err == nil {
				t.Fatalf("expected error for %s secret, got config %+v", name, cfg)
			}
			if cfg != nil {
				t.Fatalf("expected nil config on error")
			}
		})
	}
}

func TestLoadAcceptsStrongSecret(t *testing.T) {
	secret := strings.Repeat("s", 32)
	cfg, err := LoadWithGetenv(envFrom(map[string]string{"JWT_SECRET": secret}))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if cfg.JWTSecret != secret {
		t.Fatalf("secret not propagated")
	}
}

func TestLoadDefaults(t *testing.T) {
	cfg, err := LoadWithGetenv(envFrom(map[string]string{"JWT_SECRET": strings.Repeat("x", 40)}))
	if err != nil {
		t.Fatal(err)
	}
	if cfg.Port != "8086" {
		t.Errorf("Port = %q", cfg.Port)
	}
	if cfg.ListenAddr != "127.0.0.1" {
		t.Errorf("ListenAddr default = %q, want 127.0.0.1", cfg.ListenAddr)
	}
}

func TestLoadListenAddrOverride(t *testing.T) {
	cfg, err := LoadWithGetenv(envFrom(map[string]string{
		"JWT_SECRET":  strings.Repeat("x", 40),
		"LISTEN_ADDR": "0.0.0.0",
	}))
	if err != nil {
		t.Fatal(err)
	}
	if cfg.ListenAddr != "0.0.0.0" {
		t.Errorf("ListenAddr = %q", cfg.ListenAddr)
	}
}

func TestSecretErrorDoesNotLeakValue(t *testing.T) {
	secret := "short-secret-value"
	_, err := LoadWithGetenv(envFrom(map[string]string{"JWT_SECRET": secret}))
	if err == nil {
		t.Fatal("expected error")
	}
	if strings.Contains(err.Error(), secret) {
		t.Fatalf("error message leaks secret value")
	}
}
