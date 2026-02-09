package integration

import (
	"testing"

	"github.com/ulf/slka/test"
)

func TestUsersList(t *testing.T) {
	env := test.NewTestEnv(t)
	defer env.Cleanup()

	result := env.RunCommand("users", "list", "--output-pretty")
	result.AssertSuccess().AssertJSONOK()

	// Check that we got users
	users := result.GetJSONArray("data.users")
	if len(users) == 0 {
		t.Fatal("Expected some users, got none")
	}

	// Verify user structure
	first := users[0].(map[string]interface{})
	if first["id"] == nil {
		t.Error("Expected user to have 'id' field")
	}
	if first["name"] == nil {
		t.Error("Expected user to have 'name' field")
	}
}

func TestUsersListExcludeBots(t *testing.T) {
	env := test.NewTestEnv(t)
	defer env.Cleanup()

	result := env.RunCommand("users", "list", "--output-pretty")
	result.AssertSuccess().AssertJSONOK()

	users := result.GetJSONArray("data.users")

	// By default, bots should be excluded
	for _, u := range users {
		user := u.(map[string]interface{})
		if user["is_bot"] == true {
			t.Errorf("Expected bots to be excluded by default, found bot: %v", user["name"])
		}
	}
}

func TestUsersListIncludeBots(t *testing.T) {
	env := test.NewTestEnv(t)
	defer env.Cleanup()

	result := env.RunCommand("users", "list", "--include-bots", "--output-pretty")
	result.AssertSuccess().AssertJSONOK()

	users := result.GetJSONArray("data.users")

	// Should include at least one bot (Slackbot from fixtures)
	foundBot := false
	for _, u := range users {
		user := u.(map[string]interface{})
		if user["is_bot"] == true {
			foundBot = true
			break
		}
	}

	if !foundBot {
		t.Error("Expected to find at least one bot when --include-bots is used")
	}
}

func TestUsersLookupByName(t *testing.T) {
	env := test.NewTestEnv(t)
	defer env.Cleanup()

	result := env.RunCommand("users", "lookup", "alice", "--by", "name", "--output-pretty")
	result.AssertSuccess().AssertJSONOK()

	// Verify we got the right user
	user := result.GetJSONField("data.user").(map[string]interface{})
	if user["name"] != "alice" {
		t.Errorf("Expected user name to be 'alice', got %v", user["name"])
	}
	if user["id"] != "U001" {
		t.Errorf("Expected user ID to be 'U001', got %v", user["id"])
	}
}

func TestUsersLookupByEmail(t *testing.T) {
	env := test.NewTestEnv(t)
	defer env.Cleanup()

	result := env.RunCommand("users", "lookup", "alice@example.com", "--by", "email", "--output-pretty")
	result.AssertSuccess().AssertJSONOK()

	// Verify we got the right user
	user := result.GetJSONField("data.user").(map[string]interface{})
	if user["email"] != "alice@example.com" {
		t.Errorf("Expected user email to be 'alice@example.com', got %v", user["email"])
	}
}

func TestUsersLookupAuto(t *testing.T) {
	env := test.NewTestEnv(t)
	defer env.Cleanup()

	// Should work with name
	result := env.RunCommand("users", "lookup", "bob", "--output-pretty")
	result.AssertSuccess().AssertJSONOK()

	user := result.GetJSONField("data.user").(map[string]interface{})
	if user["name"] != "bob" {
		t.Errorf("Expected user name to be 'bob', got %v", user["name"])
	}
}

func TestUsersLookupNotFound(t *testing.T) {
	env := test.NewTestEnv(t)
	defer env.Cleanup()

	result := env.RunCommand("users", "lookup", "nonexistent", "--output-pretty")
	result.AssertFailure()

	// Should return an error
	data := result.ParseJSON()
	if data["ok"] != false {
		t.Error("Expected ok to be false for not found user")
	}
}
