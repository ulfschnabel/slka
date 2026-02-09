package integration

import (
	"testing"

	"github.com/ulf/slka/test"
)

func TestUnreadList(t *testing.T) {
	env := test.NewTestEnv(t)
	defer env.Cleanup()

	result := env.RunCommand("unread", "list", "--output-pretty")
	result.AssertSuccess().AssertJSONOK()

	// Check that we got unread conversations
	channels := result.GetJSONArray("data.conversations")
	if len(channels) == 0 {
		t.Fatal("Expected some unread conversations, got none")
	}

	// Verify channels are sorted by unread count (descending by default)
	// First channel should be "engineering" with 12 unread
	first := channels[0].(map[string]interface{})
	if first["name"] != "engineering" {
		t.Errorf("Expected first channel to be 'engineering', got %v", first["name"])
	}
	if first["unread_count"] != float64(12) {
		t.Errorf("Expected 12 unread messages, got %v", first["unread_count"])
	}
}

func TestUnreadListChannelsOnly(t *testing.T) {
	env := test.NewTestEnv(t)
	defer env.Cleanup()

	result := env.RunCommand("unread", "list", "--channels-only", "--output-pretty")
	result.AssertSuccess().AssertJSONOK()

	channels := result.GetJSONArray("data.conversations")

	// Verify all results are channels (not DMs)
	for _, ch := range channels {
		channel := ch.(map[string]interface{})
		if channel["is_channel"] != true {
			t.Errorf("Expected channel, got DM: %v", channel)
		}
	}
}

func TestUnreadListDMsOnly(t *testing.T) {
	env := test.NewTestEnv(t)
	defer env.Cleanup()

	result := env.RunCommand("unread", "list", "--dms-only", "--output-pretty")
	result.AssertSuccess().AssertJSONOK()

	conversations := result.GetJSONArray("data.conversations")

	// Verify all results are DMs (not channels)
	for _, conv := range conversations {
		conversation := conv.(map[string]interface{})
		isIM := conversation["is_im"] == true
		isMpIM := conversation["is_mpim"] == true
		if !isIM && !isMpIM {
			t.Errorf("Expected DM, got channel: %v", conversation)
		}
	}
}

func TestUnreadListMinUnread(t *testing.T) {
	env := test.NewTestEnv(t)
	defer env.Cleanup()

	result := env.RunCommand("unread", "list", "--min-unread", "10", "--output-pretty")
	result.AssertSuccess().AssertJSONOK()

	conversations := result.GetJSONArray("data.conversations")

	// Verify all results have at least 10 unread messages
	for _, conv := range conversations {
		conversation := conv.(map[string]interface{})
		unreadCount := int(conversation["unread_count"].(float64))
		if unreadCount < 10 {
			t.Errorf("Expected at least 10 unread, got %d", unreadCount)
		}
	}
}

func TestUnreadListOrderByOldest(t *testing.T) {
	env := test.NewTestEnv(t)
	defer env.Cleanup()

	result := env.RunCommand("unread", "list", "--order-by", "oldest", "--output-pretty")
	result.AssertSuccess().AssertJSONOK()

	conversations := result.GetJSONArray("data.conversations")
	if len(conversations) < 2 {
		t.Skip("Need at least 2 conversations to test ordering")
	}

	// Verify conversations are ordered by last_read (ascending)
	for i := 0; i < len(conversations)-1; i++ {
		current := conversations[i].(map[string]interface{})
		next := conversations[i+1].(map[string]interface{})

		currentLastRead := current["last_read"].(string)
		nextLastRead := next["last_read"].(string)

		if currentLastRead > nextLastRead {
			t.Errorf("Conversations not ordered by oldest: %s > %s", currentLastRead, nextLastRead)
		}
	}
}
