package config

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
)

// Config holds the application configuration
type Config struct {
	ReadToken       string `json:"read_token"`
	WriteToken      string `json:"write_token"`
	UserToken       string `json:"user_token,omitempty"`
	RequireApproval bool   `json:"require_approval"`
}

// Load reads configuration from file and applies environment variable overrides
func Load(path string) (*Config, error) {
	// Read config file
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}

	cfg := &Config{}
	if err := json.Unmarshal(data, cfg); err != nil {
		return nil, fmt.Errorf("failed to parse config file: %w", err)
	}

	// Apply environment variable overrides
	if token := os.Getenv("SLKA_READ_TOKEN"); token != "" {
		cfg.ReadToken = token
	}
	if token := os.Getenv("SLKA_WRITE_TOKEN"); token != "" {
		cfg.WriteToken = token
	}
	if token := os.Getenv("SLKA_USER_TOKEN"); token != "" {
		cfg.UserToken = token
	}

	// Note: require_approval cannot be overridden by environment variable

	return cfg, nil
}

// Save writes configuration to file with secure permissions
func (c *Config) Save(path string) error {
	// Create directory if it doesn't exist
	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, 0700); err != nil {
		return fmt.Errorf("failed to create config directory: %w", err)
	}

	// Marshal config to JSON
	data, err := json.MarshalIndent(c, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal config: %w", err)
	}

	// Write with secure permissions (0600 = read/write for owner only)
	if err := os.WriteFile(path, data, 0600); err != nil {
		return fmt.Errorf("failed to write config file: %w", err)
	}

	return nil
}

// Validate checks if the configuration is valid
func (c *Config) Validate() error {
	if c.ReadToken == "" {
		return errors.New("read_token is required")
	}
	if c.WriteToken == "" {
		return errors.New("write_token is required")
	}
	return nil
}

// DefaultConfigPath returns the default configuration file path
func DefaultConfigPath() string {
	home, err := os.UserHomeDir()
	if err != nil {
		return filepath.Join(".config", "slka", "config.json")
	}
	return filepath.Join(home, ".config", "slka", "config.json")
}

// MaskToken masks a token for display, showing only the prefix and last character
func MaskToken(token string) string {
	if token == "" {
		return ""
	}
	if len(token) <= 4 {
		return "***"
	}

	// Find the last dash to keep the prefix intact
	lastDash := -1
	for i := len(token) - 1; i >= 0; i-- {
		if token[i] == '-' {
			lastDash = i
			break
		}
	}

	if lastDash == -1 {
		// No dash found, just show last char
		return token[:1] + "***"
	}

	prefix := token[:lastDash+1]
	suffix := string(token[len(token)-1])
	return prefix + "***" + suffix
}
