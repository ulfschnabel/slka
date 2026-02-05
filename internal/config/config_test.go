package config

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLoadConfigFromFile(t *testing.T) {
	// Create a temporary config file
	tmpDir := t.TempDir()
	configPath := filepath.Join(tmpDir, "config.json")

	configData := `{
		"read_token": "xoxb-read-123",
		"write_token": "xoxb-write-456",
		"user_token": "xoxp-user-789",
		"require_approval": true
	}`

	err := os.WriteFile(configPath, []byte(configData), 0600)
	assert.NoError(t, err)

	cfg, err := Load(configPath)
	assert.NoError(t, err)
	assert.Equal(t, "xoxb-read-123", cfg.ReadToken)
	assert.Equal(t, "xoxb-write-456", cfg.WriteToken)
	assert.Equal(t, "xoxp-user-789", cfg.UserToken)
	assert.True(t, cfg.RequireApproval)
}

func TestLoadConfigMissingFile(t *testing.T) {
	cfg, err := Load("/nonexistent/config.json")
	assert.Error(t, err)
	assert.Nil(t, cfg)
}

func TestEnvironmentVariableOverrides(t *testing.T) {
	// Set environment variables
	os.Setenv("SLKA_READ_TOKEN", "env-read-token")
	os.Setenv("SLKA_WRITE_TOKEN", "env-write-token")
	os.Setenv("SLKA_USER_TOKEN", "env-user-token")
	defer func() {
		os.Unsetenv("SLKA_READ_TOKEN")
		os.Unsetenv("SLKA_WRITE_TOKEN")
		os.Unsetenv("SLKA_USER_TOKEN")
	}()

	// Create a temporary config file with different values
	tmpDir := t.TempDir()
	configPath := filepath.Join(tmpDir, "config.json")

	configData := `{
		"read_token": "file-read-token",
		"write_token": "file-write-token",
		"user_token": "file-user-token",
		"require_approval": true
	}`

	err := os.WriteFile(configPath, []byte(configData), 0600)
	assert.NoError(t, err)

	cfg, err := Load(configPath)
	assert.NoError(t, err)

	// Environment variables should override file values
	assert.Equal(t, "env-read-token", cfg.ReadToken)
	assert.Equal(t, "env-write-token", cfg.WriteToken)
	assert.Equal(t, "env-user-token", cfg.UserToken)
	// require_approval should still come from file (no env override)
	assert.True(t, cfg.RequireApproval)
}

func TestSaveConfig(t *testing.T) {
	tmpDir := t.TempDir()
	configPath := filepath.Join(tmpDir, "config.json")

	cfg := &Config{
		ReadToken:       "xoxb-read-new",
		WriteToken:      "xoxb-write-new",
		UserToken:       "xoxp-user-new",
		RequireApproval: false,
	}

	err := cfg.Save(configPath)
	assert.NoError(t, err)

	// Verify file was created with correct permissions
	info, err := os.Stat(configPath)
	assert.NoError(t, err)
	assert.Equal(t, os.FileMode(0600), info.Mode().Perm())

	// Load and verify
	loaded, err := Load(configPath)
	assert.NoError(t, err)
	assert.Equal(t, cfg.ReadToken, loaded.ReadToken)
	assert.Equal(t, cfg.WriteToken, loaded.WriteToken)
	assert.Equal(t, cfg.UserToken, loaded.UserToken)
	assert.Equal(t, cfg.RequireApproval, loaded.RequireApproval)
}

func TestDefaultConfigPath(t *testing.T) {
	path := DefaultConfigPath()
	assert.Contains(t, path, ".config/slka/config.json")
}

func TestMaskToken(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"xoxb-1234567890-1234567890-abcdefghijklmnop", "xoxb-***-***-***p"},
		{"xoxp-short", "xoxp-***"},
		{"", ""},
		{"x", "***"},
	}

	for _, tt := range tests {
		result := MaskToken(tt.input)
		assert.Equal(t, tt.expected, result)
	}
}

func TestConfigValidation(t *testing.T) {
	tests := []struct {
		name    string
		cfg     *Config
		wantErr bool
	}{
		{
			name: "valid config",
			cfg: &Config{
				ReadToken:  "xoxb-123",
				WriteToken: "xoxb-456",
			},
			wantErr: false,
		},
		{
			name: "missing read token",
			cfg: &Config{
				WriteToken: "xoxb-456",
			},
			wantErr: true,
		},
		{
			name: "missing write token",
			cfg: &Config{
				ReadToken: "xoxb-123",
			},
			wantErr: true,
		},
		{
			name:    "empty config",
			cfg:     &Config{},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.cfg.Validate()
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
