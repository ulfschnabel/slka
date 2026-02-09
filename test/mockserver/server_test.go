package mockserver

import (
	"encoding/json"
	"net/http"
	"testing"
)

func TestMockServerAuth(t *testing.T) {
	server := New()
	defer server.Close()

	// Test with correct token
	req, _ := http.NewRequest("GET", server.URL()+"/api/auth.test", nil)
	req.Header.Set("Authorization", "Bearer "+server.Token)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatalf("Request failed: %v", err)
	}
	defer resp.Body.Close()

	var result map[string]interface{}
	json.NewDecoder(resp.Body).Decode(&result)

	if result["ok"] != true {
		t.Errorf("Expected ok=true, got %v", result)
	}
}

func TestMockServerAuthInvalid(t *testing.T) {
	server := New()
	defer server.Close()

	// Test with wrong token
	req, _ := http.NewRequest("GET", server.URL()+"/api/auth.test", nil)
	req.Header.Set("Authorization", "Bearer wrong-token")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatalf("Request failed: %v", err)
	}
	defer resp.Body.Close()

	var result map[string]interface{}
	json.NewDecoder(resp.Body).Decode(&result)

	if result["ok"] != false {
		t.Errorf("Expected ok=false for invalid token, got %v", result)
	}
	if result["error"] != "invalid_auth" {
		t.Errorf("Expected error=invalid_auth, got %v", result["error"])
	}
}

func TestMockServerConversationsList(t *testing.T) {
	server := New()
	defer server.Close()

	req, _ := http.NewRequest("GET", server.URL()+"/api/conversations.list?types=public_channel", nil)
	req.Header.Set("Authorization", "Bearer "+server.Token)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatalf("Request failed: %v", err)
	}
	defer resp.Body.Close()

	var result map[string]interface{}
	json.NewDecoder(resp.Body).Decode(&result)

	if result["ok"] != true {
		t.Errorf("Expected ok=true, got %v. Full response: %+v", result["ok"], result)
	}

	channels := result["channels"].([]interface{})
	if len(channels) == 0 {
		t.Error("Expected some channels")
	}
}
