package approval

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestApprovalRequiredNoTTY(t *testing.T) {
	// Should fail if approval required but no TTY available
	approver := NewApprover(false, nil)

	err := approver.Require("send_message", map[string]interface{}{
		"channel": "C123",
		"text":    "Hello",
	})

	assert.Error(t, err)
	var approvalErr *ApprovalRequiredError
	assert.ErrorAs(t, err, &approvalErr)
}

func TestApprovalGrantedWithY(t *testing.T) {
	// Should proceed when user types 'y'
	mockReader := strings.NewReader("y\n")
	approver := NewApprover(true, mockReader)

	err := approver.Require("send_message", map[string]interface{}{
		"channel": "C123",
		"text":    "Hello",
	})

	assert.NoError(t, err)
}

func TestApprovalGrantedWithYes(t *testing.T) {
	// Should proceed when user types 'yes'
	mockReader := strings.NewReader("yes\n")
	approver := NewApprover(true, mockReader)

	err := approver.Require("send_message", map[string]interface{}{})

	assert.NoError(t, err)
}

func TestApprovalDenied(t *testing.T) {
	// Should abort when user types anything other than y/yes
	mockReader := strings.NewReader("n\n")
	approver := NewApprover(true, mockReader)

	err := approver.Require("send_message", map[string]interface{}{})

	assert.Error(t, err)
	var deniedErr *ApprovalDeniedError
	assert.ErrorAs(t, err, &deniedErr)
}

func TestApprovalDeniedEmpty(t *testing.T) {
	// Should abort when user just presses enter (empty input)
	mockReader := strings.NewReader("\n")
	approver := NewApprover(true, mockReader)

	err := approver.Require("send_message", map[string]interface{}{})

	assert.Error(t, err)
	var deniedErr *ApprovalDeniedError
	assert.ErrorAs(t, err, &deniedErr)
}

func TestApprovalNotRequired(t *testing.T) {
	// Should proceed without prompting when approval is not required
	approver := NewApprover(false, nil)
	approver.SetRequired(false)

	err := approver.Require("send_message", map[string]interface{}{})

	assert.NoError(t, err)
}

func TestApprovalMultipleAttempts(t *testing.T) {
	// Should keep asking until valid response
	mockReader := strings.NewReader("maybe\nno\ny\n")
	approver := NewApprover(true, mockReader)

	err := approver.Require("send_message", map[string]interface{}{})

	// After two invalid responses, third 'y' should work
	assert.NoError(t, err)
}

func TestFormatPayloadForDisplay(t *testing.T) {
	payload := map[string]interface{}{
		"channel": "C123",
		"text":    "Hello world",
		"options": map[string]interface{}{
			"unfurl_links": true,
		},
	}

	formatted := formatPayloadForDisplay(payload)
	assert.Contains(t, formatted, "C123")
	assert.Contains(t, formatted, "Hello world")
	assert.Contains(t, formatted, "unfurl_links")
}

func TestDescribeAction(t *testing.T) {
	tests := []struct {
		action  string
		payload map[string]interface{}
		want    string
	}{
		{
			action:  "send_message",
			payload: map[string]interface{}{"channel": "general"},
			want:    "Send message to general",
		},
		{
			action:  "add_reaction",
			payload: map[string]interface{}{"emoji": "thumbsup"},
			want:    "Add reaction :thumbsup:",
		},
		{
			action:  "create_channel",
			payload: map[string]interface{}{"name": "newchannel"},
			want:    "Create channel #newchannel",
		},
	}

	for _, tt := range tests {
		t.Run(tt.action, func(t *testing.T) {
			got := describeAction(tt.action, tt.payload)
			assert.Contains(t, got, tt.want)
		})
	}
}
