package integration

import (
	"testing"

	"github.com/ulf/slka/test"
)

func TestChannelsList(t *testing.T) {
	env := test.NewTestEnv(t)
	defer env.Cleanup()

	result := env.RunCommand("channels", "list", "--output-pretty")
	result.AssertSuccess().AssertJSONOK()

	channels := result.GetJSONArray("data.channels")
	if len(channels) == 0 {
		t.Fatal("Expected some channels, got none")
	}

	// Verify we have the expected test channels
	channelNames := make([]string, 0)
	for _, ch := range channels {
		channel := ch.(map[string]interface{})
		channelNames = append(channelNames, channel["name"].(string))
	}

	expectedChannels := []string{"general", "engineering", "random", "secret-project"}
	for _, expected := range expectedChannels {
		found := false
		for _, name := range channelNames {
			if name == expected {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("Expected to find channel %q, but didn't", expected)
		}
	}
}

func TestChannelsListWithFilter(t *testing.T) {
	env := test.NewTestEnv(t)
	defer env.Cleanup()

	result := env.RunCommand("channels", "list", "--filter", "eng", "--output-pretty")
	result.AssertSuccess().AssertJSONOK()

	channels := result.GetJSONArray("data.channels")

	// Should only return channels matching "eng"
	for _, ch := range channels {
		channel := ch.(map[string]interface{})
		name := channel["name"].(string)
		if name != "engineering" {
			t.Errorf("Filter 'eng' should only return 'engineering', got %q", name)
		}
	}
}

func TestChannelsInfo(t *testing.T) {
	env := test.NewTestEnv(t)
	defer env.Cleanup()

	result := env.RunCommand("channels", "info", "general", "--output-pretty")
	result.AssertSuccess().AssertJSONOK()

	// Verify channel info
	result.AssertJSONField("data.name", "general")
	result.AssertJSONField("data.id", "C001")
}

func TestChannelsHistory(t *testing.T) {
	env := test.NewTestEnv(t)
	defer env.Cleanup()

	result := env.RunCommand("channels", "history", "general", "--limit", "10", "--output-pretty")
	result.AssertSuccess().AssertJSONOK()

	messages := result.GetJSONArray("data.messages")
	if len(messages) == 0 {
		t.Error("Expected some messages in channel history")
	}

	// Verify message structure
	first := messages[0].(map[string]interface{})
	if first["text"] == nil {
		t.Error("Expected message to have 'text' field")
	}
	if first["user"] == nil {
		t.Error("Expected message to have 'user' field")
	}
}

func TestChannelsListPrivate(t *testing.T) {
	env := test.NewTestEnv(t)
	defer env.Cleanup()

	result := env.RunCommand("channels", "list", "--type", "private", "--output-pretty")
	result.AssertSuccess().AssertJSONOK()

	channels := result.GetJSONArray("data.channels")

	// Verify all channels are private
	for _, ch := range channels {
		channel := ch.(map[string]interface{})
		if channel["is_private"] != true {
			t.Errorf("Expected private channel, got public: %v", channel["name"])
		}
	}

	// Should find "secret-project"
	found := false
	for _, ch := range channels {
		channel := ch.(map[string]interface{})
		if channel["name"] == "secret-project" {
			found = true
			break
		}
	}
	if !found {
		t.Error("Expected to find 'secret-project' private channel")
	}
}
