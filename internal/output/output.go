package output

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
)

// Exit codes
const (
	ExitSuccess          = 0
	ExitGeneralError     = 1
	ExitAuthError        = 2
	ExitPermissionError  = 3
	ExitNotFound         = 4
	ExitApprovalRequired = 5
	ExitRateLimited      = 6
)

// Result represents the output structure
type Result struct {
	OK               bool                   `json:"ok"`
	Data             interface{}            `json:"data,omitempty"`
	Error            string                 `json:"error,omitempty"`
	ErrorDescription string                 `json:"error_description,omitempty"`
	Suggestion       string                 `json:"suggestion,omitempty"`
	RequiresApproval bool                   `json:"requires_approval,omitempty"`
	Action           string                 `json:"action,omitempty"`
	Description      string                 `json:"description,omitempty"`
	Payload          map[string]interface{} `json:"payload,omitempty"`
	ApproveCommand   string                 `json:"approve_command,omitempty"`
	RetryAfter       int                    `json:"retry_after,omitempty"`
	DryRun           bool                   `json:"dry_run,omitempty"`
}

// Success creates a successful result
func Success(data interface{}) *Result {
	return &Result{
		OK:   true,
		Data: data,
	}
}

// Error creates an error result
func Error(errorCode, description, suggestion string) *Result {
	return &Result{
		OK:               false,
		Error:            errorCode,
		ErrorDescription: description,
		Suggestion:       suggestion,
	}
}

// ApprovalRequired creates a result indicating approval is needed
func ApprovalRequired(action string, payload map[string]interface{}) *Result {
	return &Result{
		OK:               false,
		RequiresApproval: true,
		Action:           action,
		Description:      describeAction(action, payload),
		Payload:          payload,
	}
}

// RateLimited creates a rate limit error result
func RateLimited(retryAfter int) *Result {
	return &Result{
		OK:               false,
		Error:            "rate_limited",
		ErrorDescription: fmt.Sprintf("Rate limited by Slack API. Retry after %d seconds.", retryAfter),
		RetryAfter:       retryAfter,
	}
}

// DryRun creates a dry run result
func DryRun(action string, payload map[string]interface{}) *Result {
	return &Result{
		OK:          false,
		DryRun:      true,
		Action:      action,
		Description: describeAction(action, payload),
		Payload:     payload,
	}
}

// Format formats the result as JSON
func (r *Result) Format(pretty bool) (string, error) {
	var data []byte
	var err error

	if pretty {
		data, err = json.MarshalIndent(r, "", "  ")
	} else {
		data, err = json.Marshal(r)
	}

	if err != nil {
		return "", err
	}

	return string(data), nil
}

// Print prints the result to stdout
func (r *Result) Print(pretty bool) {
	output, err := r.Format(pretty)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error formatting output: %v\n", err)
		os.Exit(ExitGeneralError)
	}
	fmt.Println(output)
}

// ExitCode returns the appropriate exit code for this result
func (r *Result) ExitCode() int {
	if r.OK {
		return ExitSuccess
	}

	if r.RequiresApproval {
		return ExitApprovalRequired
	}

	if r.DryRun {
		return ExitSuccess
	}

	if r.RetryAfter > 0 {
		return ExitRateLimited
	}

	// Map error codes to exit codes
	switch {
	case strings.Contains(r.Error, "auth") || strings.Contains(r.Error, "token"):
		return ExitAuthError
	case strings.Contains(r.Error, "scope") || strings.Contains(r.Error, "permission"):
		return ExitPermissionError
	case strings.Contains(r.Error, "not_found") || strings.Contains(r.Error, "not_in_channel"):
		return ExitNotFound
	default:
		return ExitGeneralError
	}
}

// describeAction generates a human-readable description of an action
func describeAction(action string, payload map[string]interface{}) string {
	switch action {
	case "send_message":
		channel := payload["channel"]
		return fmt.Sprintf("Send message to %v", channel)
	case "send_dm":
		users := payload["users"]
		return fmt.Sprintf("Send DM to %v", users)
	case "add_reaction":
		emoji := payload["emoji"]
		return fmt.Sprintf("Add reaction :%v:", emoji)
	case "remove_reaction":
		emoji := payload["emoji"]
		return fmt.Sprintf("Remove reaction :%v:", emoji)
	case "create_channel":
		name := payload["name"]
		return fmt.Sprintf("Create channel #%v", name)
	case "archive_channel":
		channel := payload["channel"]
		return fmt.Sprintf("Archive channel %v", channel)
	case "invite_users":
		channel := payload["channel"]
		users := payload["users"]
		return fmt.Sprintf("Invite %v to %v", users, channel)
	default:
		return fmt.Sprintf("Execute %s", action)
	}
}
