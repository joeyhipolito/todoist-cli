package cmd

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/joeyhipolito/todoist-cli/internal/config"
)

// ConfigureCmd runs an interactive configuration setup.
func ConfigureCmd() error {
	reader := bufio.NewReader(os.Stdin)

	fmt.Println("Todoist CLI Configuration")
	fmt.Println("=========================")
	fmt.Println()

	// Check for existing config
	if config.Exists() {
		fmt.Printf("Existing configuration found at %s\n", config.Path())
		fmt.Print("Overwrite? [y/N] ")
		reply, _ := reader.ReadString('\n')
		reply = strings.TrimSpace(reply)
		if !strings.EqualFold(reply, "y") {
			fmt.Println("Configuration cancelled.")
			return nil
		}
		fmt.Println()
	}

	// Prompt for access token
	fmt.Println("Get your API token from:")
	fmt.Println("https://todoist.com/app/settings/integrations/developer")
	fmt.Println()
	fmt.Print("Todoist API Token: ")
	token, _ := reader.ReadString('\n')
	token = strings.TrimSpace(token)

	if token == "" {
		return fmt.Errorf("access token is required")
	}

	// Save configuration
	cfg := &config.Config{
		AccessToken: token,
	}

	if err := config.Save(cfg); err != nil {
		return fmt.Errorf("failed to save configuration: %w", err)
	}

	fmt.Println()
	fmt.Printf("Configuration saved to %s\n", config.Path())
	fmt.Println()
	fmt.Println("Test your setup:")
	fmt.Println("  todoist list")
	fmt.Println("  todoist projects")
	fmt.Println()
	fmt.Println("Troubleshoot:")
	fmt.Println("  todoist doctor")

	return nil
}

// ConfigureShowCmd prints the current configuration (with token masked).
func ConfigureShowCmd(jsonOutput bool) error {
	cfg, err := config.Load()
	if err != nil {
		return fmt.Errorf("failed to load config: %w", err)
	}

	if !config.Exists() {
		fmt.Println("No configuration file found.")
		fmt.Println("Run 'todoist configure' to set up.")
		return nil
	}

	// Mask token for display
	maskedToken := ""
	if cfg.AccessToken != "" {
		if len(cfg.AccessToken) > 8 {
			maskedToken = cfg.AccessToken[:4] + "..." + cfg.AccessToken[len(cfg.AccessToken)-4:]
		} else {
			maskedToken = "****"
		}
	}

	if jsonOutput {
		output := map[string]string{
			"config_path":  config.Path(),
			"access_token": maskedToken,
		}
		encoder := json.NewEncoder(os.Stdout)
		encoder.SetIndent("", "  ")
		return encoder.Encode(output)
	}

	fmt.Printf("Config file: %s\n", config.Path())
	fmt.Printf("Access token: %s\n", maskedToken)
	return nil
}
