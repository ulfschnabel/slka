package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/ulf/slka/test/fixtures"
)

type MockServer struct {
	Token    string
	channels []fixtures.Channel
	users    []fixtures.User
}

func main() {
	port := "8080"
	if p := os.Getenv("PORT"); p != "" {
		port = p
	}

	mock := &MockServer{
		Token:    "xoxp-test-token-12345",
		channels: fixtures.GetTestChannels(),
		users:    fixtures.GetTestUsers(),
	}

	mux := http.NewServeMux()
	mux.HandleFunc("/api/conversations.list", mock.handleConversationsList)
	mux.HandleFunc("/api/conversations.info", mock.handleConversationsInfo)
	mux.HandleFunc("/api/conversations.history", mock.handleConversationsHistory)
	mux.HandleFunc("/api/conversations.mark", mock.handleConversationsMark)
	mux.HandleFunc("/api/users.list", mock.handleUsersList)
	mux.HandleFunc("/api/users.info", mock.handleUsersInfo)
	mux.HandleFunc("/api/users.lookupByEmail", mock.handleUsersLookupByEmail)
	mux.HandleFunc("/api/chat.postMessage", mock.handleChatPostMessage)
	mux.HandleFunc("/api/chat.update", mock.handleChatUpdate)
	mux.HandleFunc("/api/auth.test", mock.handleAuthTest)

	fmt.Printf("🚀 Mock Slack API server starting on http://localhost:%s\n", port)
	fmt.Printf("📝 Token: %s\n", mock.Token)
	fmt.Printf("🔗 API URL: http://localhost:%s/api/\n", port)
	fmt.Println("\n✅ Ready to accept requests. Press Ctrl+C to stop.")

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)

	go func() {
		if err := http.ListenAndServe(":"+port, mux); err != nil {
			log.Fatalf("Server failed: %v", err)
		}
	}()

	<-sigChan
	fmt.Println("\n\n👋 Shutting down mock server...")
}

func (m *MockServer) checkAuth(r *http.Request) bool {
	r.ParseForm()
	token := r.FormValue("token")
	if token == "" {
		token = r.URL.Query().Get("token")
	}
	return token == m.Token
}

func (m *MockServer) writeError(w http.ResponseWriter, errMsg string) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"ok":    false,
		"error": errMsg,
	})
}

func (m *MockServer) handleConversationsList(w http.ResponseWriter, r *http.Request) {
	if !m.checkAuth(r) {
		m.writeError(w, "invalid_auth")
		return
	}

	r.ParseForm()
	types := r.FormValue("types")
	if types == "" {
		types = r.URL.Query().Get("types")
	}

	var channels []interface{}
	for _, ch := range m.channels {
		if types != "" {
			typesList := splitTypes(types)
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

		channels = append(channels, channelToAPI(ch))
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"ok":       true,
		"channels": channels,
	})
}

func (m *MockServer) handleConversationsInfo(w http.ResponseWriter, r *http.Request) {
	if !m.checkAuth(r) {
		m.writeError(w, "invalid_auth")
		return
	}

	r.ParseForm()
	channelID := r.FormValue("channel")
	if channelID == "" {
		channelID = r.URL.Query().Get("channel")
	}

	for _, ch := range m.channels {
		if ch.ID == channelID {
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(map[string]interface{}{
				"ok":      true,
				"channel": channelToAPI(ch),
			})
			return
		}
	}

	m.writeError(w, "channel_not_found")
}

func (m *MockServer) handleConversationsHistory(w http.ResponseWriter, r *http.Request) {
	if !m.checkAuth(r) {
		m.writeError(w, "invalid_auth")
		return
	}

	r.ParseForm()
	channelID := r.FormValue("channel")
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

func (m *MockServer) handleUsersList(w http.ResponseWriter, r *http.Request) {
	if !m.checkAuth(r) {
		m.writeError(w, "invalid_auth")
		return
	}

	var members []interface{}
	for _, u := range m.users {
		members = append(members, userToAPI(u))
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"ok":      true,
		"members": members,
	})
}

func (m *MockServer) handleUsersInfo(w http.ResponseWriter, r *http.Request) {
	if !m.checkAuth(r) {
		m.writeError(w, "invalid_auth")
		return
	}

	r.ParseForm()
	userID := r.FormValue("user")
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
		"user": userToAPI(*user),
	})
}

func (m *MockServer) handleUsersLookupByEmail(w http.ResponseWriter, r *http.Request) {
	if !m.checkAuth(r) {
		m.writeError(w, "invalid_auth")
		return
	}

	r.ParseForm()
	email := r.FormValue("email")
	if email == "" {
		email = r.URL.Query().Get("email")
	}

	user := fixtures.GetUserByEmail(email)
	if user == nil {
		m.writeError(w, "users_not_found")
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"ok":   true,
		"user": userToAPI(*user),
	})
}

func (m *MockServer) handleChatPostMessage(w http.ResponseWriter, r *http.Request) {
	if !m.checkAuth(r) {
		m.writeError(w, "invalid_auth")
		return
	}

	r.ParseForm()
	channel := r.FormValue("channel")
	threadTS := r.FormValue("thread_ts")

	response := map[string]interface{}{
		"ok":      true,
		"channel": channel,
		"ts":      "1234567890.000000",
	}

	if threadTS != "" {
		response["thread_ts"] = threadTS
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func (m *MockServer) handleChatUpdate(w http.ResponseWriter, r *http.Request) {
	if !m.checkAuth(r) {
		m.writeError(w, "invalid_auth")
		return
	}

	r.ParseForm()
	channel := r.FormValue("channel")
	ts := r.FormValue("ts")

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"ok":      true,
		"channel": channel,
		"ts":      ts,
	})
}

func (m *MockServer) handleAuthTest(w http.ResponseWriter, r *http.Request) {
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

func channelToAPI(ch fixtures.Channel) map[string]interface{} {
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

func userToAPI(u fixtures.User) map[string]interface{} {
	return map[string]interface{}{
		"id":        u.ID,
		"name":      u.Name,
		"real_name": u.RealName,
		"profile": map[string]interface{}{
			"email":        u.Email,
			"real_name":    u.RealName,
			"display_name": u.Name,
		},
		"is_bot": u.IsBot,
	}
}

func splitTypes(types string) []string {
	var result []string
	current := ""
	for _, ch := range types {
		if ch == ',' {
			if current != "" {
				result = append(result, current)
				current = ""
			}
		} else {
			current += string(ch)
		}
	}
	if current != "" {
		result = append(result, current)
	}
	return result
}


func (m *MockServer) handleConversationsMark(w http.ResponseWriter, r *http.Request) {
	if !m.checkAuth(r) {
		m.writeError(w, "invalid_auth")
		return
	}

	// Parse form data
	r.ParseForm()
	channelID := r.FormValue("channel")
	if channelID == "" {
		channelID = r.URL.Query().Get("channel")
	}
	ts := r.FormValue("ts")
	if ts == "" {
		ts = r.URL.Query().Get("ts")
	}

	// Validate channel exists
	found := false
	for _, ch := range m.channels {
		if ch.ID == channelID {
			found = true
			break
		}
	}

	if !found {
		m.writeError(w, "channel_not_found")
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"ok": true,
	})
}
