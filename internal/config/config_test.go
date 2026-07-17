package config

import (
	"os"
	"path/filepath"
	"testing"
)

func TestLoadUsesDotEnvValues(t *testing.T) {
	tempDir := t.TempDir()
	if err := os.WriteFile(filepath.Join(tempDir, ".env"), []byte("PORT=9090\nHOST=127.0.0.1\n"), 0o644); err != nil {
		t.Fatalf("write .env: %v", err)
	}

	oldWD, err := os.Getwd()
	if err != nil {
		t.Fatalf("getwd: %v", err)
	}
	defer func() {
		_ = os.Chdir(oldWD)
	}()

	if err := os.Chdir(tempDir); err != nil {
		t.Fatalf("chdir: %v", err)
	}

	cfg, err := Load()
	if err != nil {
		t.Fatalf("Load() error = %v", err)
	}

	if cfg.Port != "9090" {
		t.Fatalf("expected port 9090, got %q", cfg.Port)
	}
	if cfg.Host != "127.0.0.1" {
		t.Fatalf("expected host 127.0.0.1, got %q", cfg.Host)
	}
}

func TestLoadUsesDefaultsWhenEnvMissing(t *testing.T) {
	unsetEnv(t, "HOST", "PORT", "STATIC_DIR", "TEMPLATE_DIR")

	cfg, err := Load()
	if err != nil {
		t.Fatalf("Load() error = %v", err)
	}

	if cfg.Host != "0.0.0.0" {
		t.Fatalf("expected default host 0.0.0.0, got %q", cfg.Host)
	}
	if cfg.Port != "8080" {
		t.Fatalf("expected default port 8080, got %q", cfg.Port)
	}
	if cfg.StaticDir != "web/static" {
		t.Fatalf("expected default static dir web/static, got %q", cfg.StaticDir)
	}
	if cfg.TemplateDir != "web/templates" {
		t.Fatalf("expected default template dir web/templates, got %q", cfg.TemplateDir)
	}
}

func unsetEnv(t *testing.T, keys ...string) {
	t.Helper()

	for _, key := range keys {
		oldValue, hadValue := os.LookupEnv(key)
		if err := os.Unsetenv(key); err != nil {
			t.Fatalf("unset %s: %v", key, err)
		}
		t.Cleanup(func() {
			if hadValue {
				_ = os.Setenv(key, oldValue)
				return
			}
			_ = os.Unsetenv(key)
		})
	}
}
