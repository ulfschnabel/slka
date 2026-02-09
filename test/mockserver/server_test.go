package mockserver

import (
	"encoding/json"
	"net/http"
	"net/url"
	"strings"
	"testing"
)

func TestMockServerAuth(t *testing.T) {
	server := New()
	defer server.Close()

	// Test with correct token (in form data, like slack-go does)
	form := url.Values{}
	form.Add("token", server.Token)

	req, _ := http.NewRequest("POST", server.URL()+"/api/auth.test", strings.NewReader(form.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

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

	// Test with wrong token (in form data)
	form := url.Values{}
	form.Add("token", "wrong-token")

	req, _ := http.NewRequest("POST", server.URL()+"/api/auth.test", strings.NewReader(form.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

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

	// Send token and types in form data
	form := url.Values{}
	form.Add("token", server.Token)
	form.Add("types", "public_channel")

	req, _ := http.NewRequest("POST", server.URL()+"/api/conversations.list", strings.NewReader(form.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

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
