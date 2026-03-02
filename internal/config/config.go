// Package config handles reading and writing the Todoist CLI configuration file.
// Configuration is stored in ~/.todoist/config (or $TODOIST_CONFIG_DIR/config) in INI-style format.
package config

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

const (
	// ConfigDir is the default directory name for Todoist configuration.
	ConfigDir = ".todoist"
	// ConfigFile is the configuration file name.
	ConfigFile = "config"
	// EnvConfigDir is the environment variable that overrides the config directory.
	EnvConfigDir = "TODOIST_CONFIG_DIR"
)

// Config represents the Todoist CLI configuration.
type Config struct {
	AccessToken string
}

// Store holds the resolved configuration directory path.
type Store struct {
	dir string
}

// NewStoreWithEnv creates a Store using the TODOIST_CONFIG_DIR environment variable
// if set, falling back to ~/.todoist/.
func NewStoreWithEnv() (*Store, error) {
	dir := os.Getenv(EnvConfigDir)
	if dir == "" {
		home, err := os.UserHomeDir()
		if err != nil {
			return nil, fmt.Errorf("determining home directory: %w", err)
		}
		dir = filepath.Join(home, ConfigDir)
	}
	return &Store{dir: dir}, nil
}

// Dir returns the full path to the config directory.
func (s *Store) Dir() string {
	return s.dir
}

// Path returns the full path to the config file.
func (s *Store) Path() string {
	return filepath.Join(s.dir, ConfigFile)
}

// Exists returns true if the config file exists.
func (s *Store) Exists() bool {
	_, err := os.Stat(s.Path())
	return err == nil
}

// Permissions returns the file permissions of the config file, or an error.
func (s *Store) Permissions() (os.FileMode, error) {
	info, err := os.Stat(s.Path())
	if err != nil {
		return 0, err
	}
	return info.Mode().Perm(), nil
}

// Load reads the configuration from the store's config file.
// Returns an empty Config (not an error) if the file doesn't exist.
func (s *Store) Load() (*Config, error) {
	cfg := &Config{}

	f, err := os.Open(s.Path())
	if err != nil {
		if os.IsNotExist(err) {
			return cfg, nil
		}
		return nil, fmt.Errorf("opening config file: %w", err)
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())

		// Skip empty lines and comments
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		// Parse key=value
		parts := strings.SplitN(line, "=", 2)
		if len(parts) != 2 {
			continue
		}

		key := strings.TrimSpace(parts[0])
		value := strings.TrimSpace(parts[1])

		switch key {
		case "access_token":
			cfg.AccessToken = value
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("reading config file: %w", err)
	}

	return cfg, nil
}

// Save writes the configuration to the store's config file with proper permissions.
func (s *Store) Save(cfg *Config) error {
	if err := os.MkdirAll(s.dir, 0700); err != nil {
		return fmt.Errorf("creating config directory: %w", err)
	}

	var b strings.Builder
	b.WriteString("# Todoist CLI Configuration\n")
	b.WriteString("# Created by: todoist configure\n")
	b.WriteString("\n")
	b.WriteString("# Your Todoist API Token\n")
	b.WriteString("# Get from: https://todoist.com/app/settings/integrations/developer\n")
	fmt.Fprintf(&b, "access_token=%s\n", cfg.AccessToken)

	if err := os.WriteFile(s.Path(), []byte(b.String()), 0600); err != nil {
		return fmt.Errorf("writing config file: %w", err)
	}

	return nil
}

// ResolveToken returns the access token using config priority:
// config file > environment variable TODOIST_ACCESS_TOKEN.
func (s *Store) ResolveToken() string {
	cfg, err := s.Load()
	if err == nil && cfg.AccessToken != "" {
		return cfg.AccessToken
	}
	return os.Getenv("TODOIST_ACCESS_TOKEN")
}

// defaultStore returns a Store using the default env-based resolution.
// Callers that need explicit error handling should use NewStoreWithEnv directly.
func defaultStore() *Store {
	s, err := NewStoreWithEnv()
	if err != nil {
		fmt.Fprintf(os.Stderr, "config: %v\n", err)
		os.Exit(1)
	}
	return s
}

// Package-level helpers for callers that don't need a Store instance.

// Path returns the full path to the config file.
func Path() string { return defaultStore().Path() }

// Dir returns the full path to the config directory.
func Dir() string { return defaultStore().Dir() }

// Load reads the configuration from the default config file.
func Load() (*Config, error) { return defaultStore().Load() }

// Save writes the configuration to the default config file.
func Save(cfg *Config) error { return defaultStore().Save(cfg) }

// Exists returns true if the default config file exists.
func Exists() bool { return defaultStore().Exists() }

// Permissions returns the file permissions of the default config file.
func Permissions() (os.FileMode, error) { return defaultStore().Permissions() }

// ResolveToken returns the access token from config file or TODOIST_ACCESS_TOKEN env var.
func ResolveToken() string { return defaultStore().ResolveToken() }
