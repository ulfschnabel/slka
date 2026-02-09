package test

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"testing"

	"github.com/ulf/slka/test/mockserver"
)

// TestEnv sets up a complete test environment with mock server and config
type TestEnv struct {
	T          *testing.T
	MockServer *mockserver.MockSlackServer
	ConfigDir  string
	ConfigFile string
	BinaryPath string
}

// NewTestEnv creates a new test environment
func NewTestEnv(t *testing.T) *TestEnv {
	// Start mock server
	mockServer := mockserver.New()

	// Create temporary config directory
	configDir := t.TempDir()
	configFile := filepath.Join(configDir, "config.json")

	// Write test config
	config := map[string]interface{}{
		"read_token":       mockServer.Token,
		"write_token":      mockServer.Token,
		"require_approval": false,
	}
	configData, _ := json.MarshalIndent(config, "", "  ")
	if err := os.WriteFile(configFile, configData, 0600); err != nil {
		t.Fatalf("Failed to write config: %v", err)
	}

	// Build the binary if not already built
	binaryPath := filepath.Join(configDir, "slka-test")
	// Get the module root directory
	cmd := exec.Command("go", "list", "-m", "-f", "{{.Dir}}")
	output, err := cmd.Output()
	if err != nil {
		t.Fatalf("Failed to find module root: %v", err)
	}
	moduleRoot := string(output[:len(output)-1]) // trim newline

	buildCmd := exec.Command("go", "build", "-o", binaryPath, filepath.Join(moduleRoot, "cmd/slka"))
	if err := buildCmd.Run(); err != nil {
		t.Fatalf("Failed to build slka: %v", err)
	}

	// Override Slack API URL using environment variable
	// (Note: The slack-go library needs to be configured to use custom URLs)

	return &TestEnv{
		T:          t,
		MockServer: mockServer,
		ConfigDir:  configDir,
		ConfigFile: configFile,
		BinaryPath: binaryPath,
	}
}

// Cleanup tears down the test environment
func (e *TestEnv) Cleanup() {
	e.MockServer.Close()
}

// RunCommand executes a slka command and returns the output
func (e *TestEnv) RunCommand(args ...string) *CommandResult {
	// Add config flag
	fullArgs := append([]string{"--config", e.ConfigFile}, args...)

	cmd := exec.Command(e.BinaryPath, fullArgs...)

	// Set SLACK_API_URL environment variable to point to mock server
	// Note: slack-go expects the base URL with /api/ suffix
	apiURL := e.MockServer.URL() + "/api/"
	cmd.Env = append(os.Environ(),
		fmt.Sprintf("SLACK_API_URL=%s", apiURL),
		"SLKA_DEBUG=1",
	)

	// Debug: print the API URL being used
	e.T.Logf("Using SLACK_API_URL: %s", apiURL)
	e.T.Logf("Using token: %s", e.MockServer.Token)

	// Capture stdout and stderr separately
	// JSON output goes to stdout, debug output goes to stderr
	output, err := cmd.Output() // Only capture stdout

	return &CommandResult{
		T:        e.T,
		Output:   string(output),
		ExitCode: cmd.ProcessState.ExitCode(),
		Err:      err,
	}
}

// CommandResult represents the result of running a CLI command
type CommandResult struct {
	T        *testing.T
	Output   string
	ExitCode int
	Err      error
}

// AssertSuccess asserts the command succeeded
func (r *CommandResult) AssertSuccess() *CommandResult {
	if r.ExitCode != 0 {
		r.T.Fatalf("Command failed with exit code %d: %s", r.ExitCode, r.Output)
	}
	return r
}

// AssertFailure asserts the command failed
func (r *CommandResult) AssertFailure() *CommandResult {
	if r.ExitCode == 0 {
		r.T.Fatal("Expected command to fail, but it succeeded")
	}
	return r
}

// AssertContains asserts the output contains a substring
func (r *CommandResult) AssertContains(substr string) *CommandResult {
	if !contains(r.Output, substr) {
		r.T.Fatalf("Output does not contain %q:\n%s", substr, r.Output)
	}
	return r
}

// AssertNotContains asserts the output does not contain a substring
func (r *CommandResult) AssertNotContains(substr string) *CommandResult {
	if contains(r.Output, substr) {
		r.T.Fatalf("Output should not contain %q:\n%s", substr, r.Output)
	}
	return r
}

// ParseJSON parses the output as JSON
func (r *CommandResult) ParseJSON() map[string]interface{} {
	var data map[string]interface{}
	if err := json.Unmarshal([]byte(r.Output), &data); err != nil {
		r.T.Fatalf("Failed to parse JSON output: %v\nOutput: %s", err, r.Output)
	}
	return data
}

// AssertJSONField asserts a JSON field has a specific value
func (r *CommandResult) AssertJSONField(path string, expected interface{}) *CommandResult {
	data := r.ParseJSON()

	// Simple path navigation (supports dot notation like "data.channels")
	keys := splitPath(path)
	current := data

	for i, key := range keys {
		if i == len(keys)-1 {
			// Last key - check value
			if current[key] != expected {
				r.T.Fatalf("Expected %s to be %v, got %v", path, expected, current[key])
			}
		} else {
			// Navigate deeper
			next, ok := current[key].(map[string]interface{})
			if !ok {
				r.T.Fatalf("Cannot navigate to %s in JSON", path)
			}
			current = next
		}
	}

	return r
}

// AssertJSONOK asserts the JSON has ok=true
func (r *CommandResult) AssertJSONOK() *CommandResult {
	return r.AssertJSONField("ok", true)
}

// GetJSONField gets a field from the JSON output
func (r *CommandResult) GetJSONField(path string) interface{} {
	data := r.ParseJSON()
	keys := splitPath(path)
	current := data

	for i, key := range keys {
		if i == len(keys)-1 {
			return current[key]
		}
		next, ok := current[key].(map[string]interface{})
		if !ok {
			r.T.Fatalf("Cannot navigate to %s in JSON", path)
		}
		current = next
	}

	return nil
}

// GetJSONArray gets an array field from the JSON output
func (r *CommandResult) GetJSONArray(path string) []interface{} {
	value := r.GetJSONField(path)
	arr, ok := value.([]interface{})
	if !ok {
		r.T.Fatalf("Field %s is not an array. Value: %+v (type: %T). Full output: %s", path, value, value, r.Output)
	}
	return arr
}

// Helper functions

func contains(s, substr string) bool {
	return len(s) >= len(substr) &&
		(s == substr || len(s) > len(substr) && stringContains(s, substr))
}

func stringContains(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}

func splitPath(path string) []string {
	var result []string
	current := ""
	for _, ch := range path {
		if ch == '.' {
			if current != "" {
				result = append(result, current)
				current = ""
			}
		} else {
			current += string(ch)
		}
	}
	if current != "" {
		result = append(result, current)
	}
	return result
}
