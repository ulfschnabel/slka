package approval

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"strings"
)

// ApprovalRequiredError is returned when approval is required but not available
type ApprovalRequiredError struct {
	Action  string
	Payload map[string]interface{}
}

func (e *ApprovalRequiredError) Error() string {
	return fmt.Sprintf("approval required for action: %s", e.Action)
}

// ApprovalDeniedError is returned when user explicitly denies approval
type ApprovalDeniedError struct {
	Action string
}

func (e *ApprovalDeniedError) Error() string {
	return fmt.Sprintf("approval denied for action: %s", e.Action)
}

// Approver handles human approval workflow
type Approver struct {
	isatty   bool
	reader   io.Reader
	required bool
}

// NewApprover creates a new approver
func NewApprover(isatty bool, reader io.Reader) *Approver {
	if reader == nil {
		reader = os.Stdin
	}
	return &Approver{
		isatty:   isatty,
		reader:   reader,
		required: true, // default to requiring approval
	}
}

// SetRequired sets whether approval is required
func (a *Approver) SetRequired(required bool) {
	a.required = required
}

// Require prompts for approval if needed
func (a *Approver) Require(action string, payload map[string]interface{}) error {
	// If approval is not required, proceed immediately
	if !a.required {
		return nil
	}

	// If no TTY available, return error
	if !a.isatty {
		return &ApprovalRequiredError{
			Action:  action,
			Payload: payload,
		}
	}

	// Display what will be done
	fmt.Fprintf(os.Stderr, "\n%s\n", describeAction(action, payload))
	fmt.Fprintf(os.Stderr, "\nPayload:\n%s\n", formatPayloadForDisplay(payload))
	fmt.Fprintf(os.Stderr, "\nExecute this action? [y/N]: ")

	// Read user input
	scanner := bufio.NewScanner(a.reader)
	for scanner.Scan() {
		response := strings.TrimSpace(strings.ToLower(scanner.Text()))

		if response == "y" || response == "yes" {
			return nil
		}

		if response == "n" || response == "no" || response == "" {
			return &ApprovalDeniedError{Action: action}
		}

		// Invalid response, ask again
		fmt.Fprintf(os.Stderr, "Please enter 'y' or 'n': ")
	}

	if err := scanner.Err(); err != nil {
		return fmt.Errorf("error reading input: %w", err)
	}

	return &ApprovalDeniedError{Action: action}
}

// describeAction generates a human-readable description
func describeAction(action string, payload map[string]interface{}) string {
	switch action {
	case "send_message":
		channel := payload["channel"]
		text := payload["text"]
		return fmt.Sprintf("Send message to %v: %q", channel, text)
	case "send_dm":
		users := payload["users"]
		text := payload["text"]
		return fmt.Sprintf("Send DM to %v: %q", users, text)
	case "add_reaction":
		emoji := payload["emoji"]
		channel := payload["channel"]
		return fmt.Sprintf("Add reaction :%v: to message in %v", emoji, channel)
	case "remove_reaction":
		emoji := payload["emoji"]
		channel := payload["channel"]
		return fmt.Sprintf("Remove reaction :%v: from message in %v", emoji, channel)
	case "create_channel":
		name := payload["name"]
		private := payload["private"]
		if private == true {
			return fmt.Sprintf("Create private channel #%v", name)
		}
		return fmt.Sprintf("Create channel #%v", name)
	case "archive_channel":
		channel := payload["channel"]
		return fmt.Sprintf("Archive channel %v", channel)
	case "unarchive_channel":
		channel := payload["channel"]
		return fmt.Sprintf("Unarchive channel %v", channel)
	case "rename_channel":
		channel := payload["channel"]
		newName := payload["new_name"]
		return fmt.Sprintf("Rename channel %v to %v", channel, newName)
	case "invite_users":
		channel := payload["channel"]
		users := payload["users"]
		return fmt.Sprintf("Invite %v to %v", users, channel)
	case "kick_users":
		channel := payload["channel"]
		users := payload["users"]
		return fmt.Sprintf("Remove %v from %v", users, channel)
	default:
		return fmt.Sprintf("Execute action: %s", action)
	}
}

// formatPayloadForDisplay formats the payload as pretty JSON
func formatPayloadForDisplay(payload map[string]interface{}) string {
	data, err := json.MarshalIndent(payload, "", "  ")
	if err != nil {
		return fmt.Sprintf("%v", payload)
	}
	return string(data)
}
