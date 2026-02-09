package integration

import (
	"testing"

	"github.com/ulf/slka/test"
)

func TestMessageSend(t *testing.T) {
	env := test.NewTestEnv(t)
	defer env.Cleanup()

	result := env.RunCommand("message", "send", "general", "Hello from test!", "--output-pretty")
	result.AssertSuccess().AssertJSONOK()

	// Verify response includes channel and timestamp
	channel := result.GetJSONField("data.channel")
	if channel == nil {
		t.Error("Expected channel in response")
	}

	ts := result.GetJSONField("data.ts")
	if ts == nil {
		t.Error("Expected ts (timestamp) in response")
	}
}

func TestMessageSendByChannelID(t *testing.T) {
	env := test.NewTestEnv(t)
	defer env.Cleanup()

	// Should work with channel ID instead of name
	result := env.RunCommand("message", "send", "C001", "Test message", "--output-pretty")
	result.AssertSuccess().AssertJSONOK()

	channel := result.GetJSONField("data.channel")
	if channel != "C001" {
		t.Errorf("Expected channel C001, got %v", channel)
	}
}

func TestMessageSendInvalidChannel(t *testing.T) {
	env := test.NewTestEnv(t)
	defer env.Cleanup()

	result := env.RunCommand("message", "send", "nonexistent", "Test", "--output-pretty")
	result.AssertFailure()

	data := result.ParseJSON()
	if data["ok"] != false {
		t.Error("Expected ok to be false for invalid channel")
	}
}

func TestMessageSendWithUnfurlFlags(t *testing.T) {
	env := test.NewTestEnv(t)
	defer env.Cleanup()

	result := env.RunCommand("message", "send", "general",
		"Check https://example.com",
		"--unfurl-links=false",
		"--unfurl-media=false",
		"--output-pretty")
	result.AssertSuccess().AssertJSONOK()
}

func TestMessageReply(t *testing.T) {
	env := test.NewTestEnv(t)
	defer env.Cleanup()

	result := env.RunCommand("message", "reply", "general", "1234567890.000000", "Reply message", "--output-pretty")
	result.AssertSuccess().AssertJSONOK()

	// Verify thread_ts is returned
	threadTS := result.GetJSONField("data.thread_ts")
	if threadTS != "1234567890.000000" {
		t.Errorf("Expected thread_ts to be 1234567890.000000, got %v", threadTS)
	}
}

func TestMessageEdit(t *testing.T) {
	env := test.NewTestEnv(t)
	defer env.Cleanup()

	result := env.RunCommand("message", "edit", "general", "1234567890.000000", "Edited message", "--output-pretty")
	result.AssertSuccess().AssertJSONOK()

	// Verify timestamp is returned
	ts := result.GetJSONField("data.ts")
	if ts == nil {
		t.Error("Expected ts in response")
	}
}
