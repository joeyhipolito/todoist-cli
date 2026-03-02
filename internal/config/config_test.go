package config

import (
	"os"
	"path/filepath"
	"testing"
)

func TestNewStoreWithEnv_default(t *testing.T) {
	// Ensure env var is unset.
	os.Unsetenv(EnvConfigDir)

	s, err := NewStoreWithEnv()
	if err != nil {
		t.Fatalf("NewStoreWithEnv() error = %v", err)
	}

	home, err := os.UserHomeDir()
	if err != nil {
		t.Fatalf("UserHomeDir() error = %v", err)
	}
	want := filepath.Join(home, ConfigDir)
	if s.Dir() != want {
		t.Errorf("Dir() = %q, want %q", s.Dir(), want)
	}
	if s.Path() != filepath.Join(want, ConfigFile) {
		t.Errorf("Path() = %q, want %q", s.Path(), filepath.Join(want, ConfigFile))
	}
}

func TestNewStoreWithEnv_envOverride(t *testing.T) {
	tmp := t.TempDir()
	t.Setenv(EnvConfigDir, tmp)

	s, err := NewStoreWithEnv()
	if err != nil {
		t.Fatalf("NewStoreWithEnv() error = %v", err)
	}

	if s.Dir() != tmp {
		t.Errorf("Dir() = %q, want %q", s.Dir(), tmp)
	}
	if s.Path() != filepath.Join(tmp, ConfigFile) {
		t.Errorf("Path() = %q, want %q", s.Path(), filepath.Join(tmp, ConfigFile))
	}
}

func TestStore_LoadSave(t *testing.T) {
	tmp := t.TempDir()
	t.Setenv(EnvConfigDir, tmp)

	s, err := NewStoreWithEnv()
	if err != nil {
		t.Fatalf("NewStoreWithEnv() error = %v", err)
	}

	// Load from non-existent file returns empty config (no error).
	cfg, err := s.Load()
	if err != nil {
		t.Fatalf("Load() on missing file error = %v", err)
	}
	if cfg.AccessToken != "" {
		t.Errorf("Load() AccessToken = %q, want empty", cfg.AccessToken)
	}

	// Save and reload.
	want := &Config{AccessToken: "test-token-abc"}
	if err := s.Save(want); err != nil {
		t.Fatalf("Save() error = %v", err)
	}

	got, err := s.Load()
	if err != nil {
		t.Fatalf("Load() after Save() error = %v", err)
	}
	if got.AccessToken != want.AccessToken {
		t.Errorf("Load() AccessToken = %q, want %q", got.AccessToken, want.AccessToken)
	}
}

func TestStore_Exists(t *testing.T) {
	tmp := t.TempDir()
	t.Setenv(EnvConfigDir, tmp)

	s, err := NewStoreWithEnv()
	if err != nil {
		t.Fatalf("NewStoreWithEnv() error = %v", err)
	}

	if s.Exists() {
		t.Error("Exists() = true before any config written")
	}

	if err := s.Save(&Config{AccessToken: "tok"}); err != nil {
		t.Fatalf("Save() error = %v", err)
	}

	if !s.Exists() {
		t.Error("Exists() = false after config written")
	}
}

func TestStore_ResolveToken(t *testing.T) {
	tmp := t.TempDir()
	t.Setenv(EnvConfigDir, tmp)

	s, err := NewStoreWithEnv()
	if err != nil {
		t.Fatalf("NewStoreWithEnv() error = %v", err)
	}

	// No config, no env — empty.
	os.Unsetenv("TODOIST_ACCESS_TOKEN")
	if tok := s.ResolveToken(); tok != "" {
		t.Errorf("ResolveToken() = %q, want empty", tok)
	}

	// Env var fallback.
	t.Setenv("TODOIST_ACCESS_TOKEN", "env-token")
	if tok := s.ResolveToken(); tok != "env-token" {
		t.Errorf("ResolveToken() = %q, want %q", tok, "env-token")
	}

	// Config file takes priority.
	if err := s.Save(&Config{AccessToken: "file-token"}); err != nil {
		t.Fatalf("Save() error = %v", err)
	}
	if tok := s.ResolveToken(); tok != "file-token" {
		t.Errorf("ResolveToken() = %q, want %q", tok, "file-token")
	}
}
