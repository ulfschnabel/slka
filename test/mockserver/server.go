package mockserver

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"

	"github.com/ulf/slka/test/fixtures"
)

// MockSlackServer simulates the Slack API for testing
type MockSlackServer struct {
	Server   *httptest.Server
	Token    string
	channels []fixtures.Channel
	users    []fixtures.User
}

// New creates a new mock Slack API server
func New() *MockSlackServer {
	m := &MockSlackServer{
		Token:    "xoxp-test-token-12345",
		channels: fixtures.GetTestChannels(),
		users:    fixtures.GetTestUsers(),
	}

	mux := http.NewServeMux()
	mux.HandleFunc("/api/conversations.list", m.handleConversationsList)
	mux.HandleFunc("/api/conversations.info", m.handleConversationsInfo)
	mux.HandleFunc("/api/conversations.history", m.handleConversationsHistory)
	mux.HandleFunc("/api/users.list", m.handleUsersList)
	mux.HandleFunc("/api/users.info", m.handleUsersInfo)
	mux.HandleFunc("/api/chat.postMessage", m.handleChatPostMessage)
	mux.HandleFunc("/api/auth.test", m.handleAuthTest)

	m.Server = httptest.NewServer(mux)
	return m
}

// Close shuts down the mock server
func (m *MockSlackServer) Close() {
	m.Server.Close()
}

// URL returns the base URL of the mock server
func (m *MockSlackServer) URL() string {
	return m.Server.URL
}

// checkAuth validates the token from form data or query parameter
func (m *MockSlackServer) checkAuth(r *http.Request) bool {
	// Parse form data to get token
	r.ParseForm()
	token := r.FormValue("token")

	// Also check query parameter as fallback
	if token == "" {
		token = r.URL.Query().Get("token")
	}

	return token == m.Token
}

// writeError writes an error response
func (m *MockSlackServer) writeError(w http.ResponseWriter, errMsg string) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"ok":    false,
		"error": errMsg,
	})
}

// handleConversationsList handles conversations.list API calls
func (m *MockSlackServer) handleConversationsList(w http.ResponseWriter, r *http.Request) {
	if !m.checkAuth(r) {
		m.writeError(w, "invalid_auth")
		return
	}

	// Parse form data to get types
	r.ParseForm()
	types := r.FormValue("types")
	// Also check query parameter as fallback
	if types == "" {
		types = r.URL.Query().Get("types")
	}

	var channels []interface{}
	for _, ch := range m.channels {
		// Filter by types
		if types != "" {
			typesList := strings.Split(types, ",")
			matched := false
			for _, t := range typesList {
				if (t == "public_channel" && ch.IsChannel && !ch.IsPrivate) ||
					(t == "private_channel" && ch.IsChannel && ch.IsPrivate) ||
					(t == "im" && ch.IsIM) ||
					(t == "mpim" && ch.IsMpIM) {
					matched = true
					break
				}
			}
			if !matched {
				continue
			}
		}

		channels = append(channels, m.channelToAPI(ch))
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"ok":       true,
		"channels": channels,
	})
}

// handleConversationsInfo handles conversations.info API calls
func (m *MockSlackServer) handleConversationsInfo(w http.ResponseWriter, r *http.Request) {
	if !m.checkAuth(r) {
		m.writeError(w, "invalid_auth")
		return
	}

	// Parse form data to get channel ID
	r.ParseForm()
	channelID := r.FormValue("channel")
	// Also check query parameter as fallback
	if channelID == "" {
		channelID = r.URL.Query().Get("channel")
	}
	if channelID == "" {
		m.writeError(w, "channel_not_found")
		return
	}

	for _, ch := range m.channels {
		if ch.ID == channelID {
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(map[string]interface{}{
				"ok":      true,
				"channel": m.channelToAPI(ch),
			})
			return
		}
	}

	m.writeError(w, "channel_not_found")
}

// handleConversationsHistory handles conversations.history API calls
func (m *MockSlackServer) handleConversationsHistory(w http.ResponseWriter, r *http.Request) {
	if !m.checkAuth(r) {
		m.writeError(w, "invalid_auth")
		return
	}

	// Parse form data to get channel ID
	r.ParseForm()
	channelID := r.FormValue("channel")
	// Also check query parameter as fallback
	if channelID == "" {
		channelID = r.URL.Query().Get("channel")
	}
	messages := fixtures.GetTestMessages(channelID)

	var apiMessages []interface{}
	for _, msg := range messages {
		apiMessages = append(apiMessages, map[string]interface{}{
			"type": msg.Type,
			"user": msg.User,
			"text": msg.Text,
			"ts":   msg.Timestamp,
		})
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"ok":       true,
		"messages": apiMessages,
	})
}

// handleUsersList handles users.list API calls
func (m *MockSlackServer) handleUsersList(w http.ResponseWriter, r *http.Request) {
	if !m.checkAuth(r) {
		m.writeError(w, "invalid_auth")
		return
	}

	var members []interface{}
	for _, u := range m.users {
		members = append(members, m.userToAPI(u))
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"ok":      true,
		"members": members,
	})
}

// handleUsersInfo handles users.info API calls
func (m *MockSlackServer) handleUsersInfo(w http.ResponseWriter, r *http.Request) {
	if !m.checkAuth(r) {
		m.writeError(w, "invalid_auth")
		return
	}

	// Parse form data to get user ID
	r.ParseForm()
	userID := r.FormValue("user")
	// Also check query parameter as fallback
	if userID == "" {
		userID = r.URL.Query().Get("user")
	}
	user := fixtures.GetUserByID(userID)
	if user == nil {
		m.writeError(w, "user_not_found")
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"ok":   true,
		"user": m.userToAPI(*user),
	})
}

// handleChatPostMessage handles chat.postMessage API calls
func (m *MockSlackServer) handleChatPostMessage(w http.ResponseWriter, r *http.Request) {
	if !m.checkAuth(r) {
		m.writeError(w, "invalid_auth")
		return
	}

	// Just return success for now
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"ok":      true,
		"channel": "C001",
		"ts":      "1234567890.000000",
	})
}

// handleAuthTest handles auth.test API calls
func (m *MockSlackServer) handleAuthTest(w http.ResponseWriter, r *http.Request) {
	if !m.checkAuth(r) {
		m.writeError(w, "invalid_auth")
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"ok":      true,
		"url":     "https://testworkspace.slack.com/",
		"team":    "Test Workspace",
		"user":    "testuser",
		"team_id": "T12345",
		"user_id": "U12345",
	})
}

// channelToAPI converts a fixture channel to Slack API format
func (m *MockSlackServer) channelToAPI(ch fixtures.Channel) map[string]interface{}{
	return map[string]interface{}{
		"id":                   ch.ID,
		"name":                 ch.Name,
		"is_channel":           ch.IsChannel,
		"is_private":           ch.IsPrivate,
		"is_im":                ch.IsIM,
		"is_mpim":              ch.IsMpIM,
		"unread_count":         ch.UnreadCount,
		"unread_count_display": ch.UnreadCountDisplay,
		"last_read":            ch.LastRead,
		"user":                 ch.User,
		"num_members":          ch.NumMembers,
	}
}

// userToAPI converts a fixture user to Slack API format
func (m *MockSlackServer) userToAPI(u fixtures.User) map[string]interface{}{
	return map[string]interface{}{
		"id":        u.ID,
		"name":      u.Name,
		"real_name": u.RealName,
		"profile": map[string]interface{}{
			"email":      u.Email,
			"real_name":  u.RealName,
			"display_name": u.Name,
		},
		"is_bot": u.IsBot,
	}
}
