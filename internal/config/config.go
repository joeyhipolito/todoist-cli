// Package config handles reading and writing the Todoist CLI configuration file.
// Configuration is stored in ~/.todoist/config in INI-style format.
package config

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

const (
	// ConfigDir is the directory name for Todoist configuration.
	ConfigDir = ".todoist"
	// ConfigFile is the configuration file name.
	ConfigFile = "config"
)

// Config represents the Todoist CLI configuration.
type Config struct {
	AccessToken string
}

// Path returns the full path to the config file (~/.todoist/config).
func Path() string {
	home, err := os.UserHomeDir()
	if err != nil {
		return ""
	}
	return filepath.Join(home, ConfigDir, ConfigFile)
}

// Dir returns the full path to the config directory (~/.todoist/).
func Dir() string {
	home, err := os.UserHomeDir()
	if err != nil {
		return ""
	}
	return filepath.Join(home, ConfigDir)
}

// Load reads the configuration from ~/.todoist/config.
// Returns an empty Config (not an error) if the file doesn't exist.
func Load() (*Config, error) {
	cfg := &Config{}
	path := Path()
	if path == "" {
		return cfg, nil
	}

	f, err := os.Open(path)
	if err != nil {
		if os.IsNotExist(err) {
			return cfg, nil
		}
		return nil, fmt.Errorf("failed to open config file: %w", err)
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
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}

	return cfg, nil
}

// Save writes the configuration to ~/.todoist/config with proper permissions.
func Save(cfg *Config) error {
	dir := Dir()
	if dir == "" {
		return fmt.Errorf("cannot determine home directory")
	}

	// Create config directory with 700 permissions
	if err := os.MkdirAll(dir, 0700); err != nil {
		return fmt.Errorf("failed to create config directory: %w", err)
	}

	path := Path()

	// Build config content
	var b strings.Builder
	b.WriteString("# Todoist CLI Configuration\n")
	b.WriteString("# Created by: todoist configure\n")
	b.WriteString("\n")
	b.WriteString("# Your Todoist API Token\n")
	b.WriteString("# Get from: https://todoist.com/app/settings/integrations/developer\n")
	fmt.Fprintf(&b, "access_token=%s\n", cfg.AccessToken)

	// Write file with 600 permissions
	if err := os.WriteFile(path, []byte(b.String()), 0600); err != nil {
		return fmt.Errorf("failed to write config file: %w", err)
	}

	return nil
}

// Exists returns true if the config file exists.
func Exists() bool {
	path := Path()
	if path == "" {
		return false
	}
	_, err := os.Stat(path)
	return err == nil
}

// Permissions returns the file permissions of the config file, or an error.
func Permissions() (os.FileMode, error) {
	path := Path()
	if path == "" {
		return 0, fmt.Errorf("cannot determine config path")
	}
	info, err := os.Stat(path)
	if err != nil {
		return 0, err
	}
	return info.Mode().Perm(), nil
}

// ResolveToken returns the access token using config priority:
// config file > environment variable.
func ResolveToken() string {
	cfg, err := Load()
	if err == nil && cfg.AccessToken != "" {
		return cfg.AccessToken
	}
	return os.Getenv("TODOIST_ACCESS_TOKEN")
}
