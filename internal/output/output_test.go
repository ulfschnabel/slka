package output

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSuccessResponse(t *testing.T) {
	result := Success(map[string]interface{}{
		"message": "test",
	})

	assert.True(t, result.OK)
	assert.Equal(t, "test", result.Data.(map[string]interface{})["message"])
}

func TestErrorResponse(t *testing.T) {
	result := Error("test_error", "Test error description", "Try again")

	assert.False(t, result.OK)
	assert.Equal(t, "test_error", result.Error)
	assert.Equal(t, "Test error description", result.ErrorDescription)
	assert.Equal(t, "Try again", result.Suggestion)
}

func TestFormatJSONCompact(t *testing.T) {
	result := Success(map[string]interface{}{
		"test": "value",
	})

	json, err := result.Format(false)
	assert.NoError(t, err)
	assert.Contains(t, json, `"ok":true`)
	assert.Contains(t, json, `"test":"value"`)
	assert.NotContains(t, json, "\n")
}

func TestFormatJSONPretty(t *testing.T) {
	result := Success(map[string]interface{}{
		"test": "value",
	})

	json, err := result.Format(true)
	assert.NoError(t, err)
	assert.Contains(t, json, `"ok": true`)
	assert.Contains(t, json, `"test": "value"`)
	assert.Contains(t, json, "\n")
}

func TestPrintSuccess(t *testing.T) {
	result := Success(map[string]interface{}{
		"test": "value",
	})

	// This should not panic
	assert.NotPanics(t, func() {
		result.Print(false)
	})
}

func TestExitCodeSuccess(t *testing.T) {
	result := Success(nil)
	assert.Equal(t, ExitSuccess, result.ExitCode())
}

func TestExitCodeError(t *testing.T) {
	result := Error("general_error", "test", "")
	assert.Equal(t, ExitGeneralError, result.ExitCode())
}

func TestExitCodeAuthError(t *testing.T) {
	result := Error("invalid_auth", "test", "")
	assert.Equal(t, ExitAuthError, result.ExitCode())
}

func TestExitCodePermissionError(t *testing.T) {
	result := Error("missing_scope", "test", "")
	assert.Equal(t, ExitPermissionError, result.ExitCode())
}

func TestExitCodeNotFound(t *testing.T) {
	result := Error("channel_not_found", "test", "")
	assert.Equal(t, ExitNotFound, result.ExitCode())
}

func TestExitCodeApprovalRequired(t *testing.T) {
	result := ApprovalRequired("test_action", map[string]interface{}{})
	assert.Equal(t, ExitApprovalRequired, result.ExitCode())
}

func TestExitCodeRateLimited(t *testing.T) {
	result := RateLimited(60)
	assert.Equal(t, ExitRateLimited, result.ExitCode())
}

func TestApprovalRequiredResponse(t *testing.T) {
	payload := map[string]interface{}{
		"channel": "C123",
		"text":    "Hello",
	}
	result := ApprovalRequired("send_message", payload)

	assert.False(t, result.OK)
	assert.True(t, result.RequiresApproval)
	assert.Equal(t, "send_message", result.Action)
	assert.NotEmpty(t, result.Description)
	assert.Equal(t, payload, result.Payload)
}

func TestRateLimitedResponse(t *testing.T) {
	result := RateLimited(120)

	assert.False(t, result.OK)
	assert.Equal(t, "rate_limited", result.Error)
	assert.Equal(t, 120, result.RetryAfter)
}

func TestDryRunResponse(t *testing.T) {
	payload := map[string]interface{}{
		"channel": "C123",
	}
	result := DryRun("send_message", payload)

	assert.False(t, result.OK)
	assert.True(t, result.DryRun)
	assert.Equal(t, "send_message", result.Action)
	assert.Equal(t, payload, result.Payload)
}
