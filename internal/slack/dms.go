package slack

import (
	"fmt"
	"strings"

	"github.com/slack-go/slack"
)

// DMInfo represents a direct message conversation
type DMInfo struct {
	ID       string `json:"id"`
	UserID   string `json:"user_id"`
	UserName string `json:"user_name,omitempty"`
}

// DMService provides DM operations
type DMService struct {
	client Client
}

// NewDMService creates a new DMService
func NewDMService(client Client) *DMService {
	return &DMService{client: client}
}

// List returns all DM conversations
func (s *DMService) List(limit int) ([]DMInfo, error) {
	params := &slack.GetConversationsParameters{
		Types: []string{"im"},
		Limit: limit,
	}

	channels, _, err := s.client.GetConversations(params)
	if err != nil {
		return nil, fmt.Errorf("failed to list DMs: %w", err)
	}

	dms := make([]DMInfo, 0, len(channels))
	for _, channel := range channels {
		dm := DMInfo{
			ID:     channel.ID,
			UserID: channel.User,
		}

		// Try to get user name
		if channel.User != "" {
			if user, err := s.client.GetUserInfo(channel.User); err == nil {
				dm.UserName = user.Name
			}
		}

		dms = append(dms, dm)
	}

	return dms, nil
}

// ResolveUser converts a user identifier (ID, email, or name) to a user ID
func (s *DMService) ResolveUser(userArg string) (string, error) {
	// If it's already a user ID (starts with U), return it
	if strings.HasPrefix(userArg, "U") {
		return userArg, nil
	}

	// Try as email first (contains @)
	if strings.Contains(userArg, "@") {
		user, err := s.client.GetUserByEmail(userArg)
		if err == nil {
			return user.ID, nil
		}
	}

	// Try as email even without @ (in case it's an email without domain shown)
	user, err := s.client.GetUserByEmail(userArg)
	if err == nil {
		return user.ID, nil
	}

	// Try as username - get all users and find by name
	users, err := s.client.GetUsers()
	if err != nil {
		return "", fmt.Errorf("failed to get users: %w", err)
	}

	for _, u := range users {
		if u.Name == userArg {
			return u.ID, nil
		}
	}

	return "", fmt.Errorf("user not found: %s", userArg)
}

// GetHistory gets the message history of a DM conversation with a user
func (s *DMService) GetHistory(userArg string, options HistoryOptions) ([]MessageInfo, error) {
	// Resolve user to ID
	userID, err := s.ResolveUser(userArg)
	if err != nil {
		return nil, err
	}

	// Open/get the DM conversation
	conv, _, _, err := s.client.OpenConversation(&slack.OpenConversationParameters{
		Users: []string{userID},
	})
	if err != nil {
		return nil, fmt.Errorf("failed to open DM with user: %w", err)
	}

	// Get conversation history using existing channel service patterns
	params := &slack.GetConversationHistoryParameters{
		ChannelID: conv.ID,
		Limit:     options.Limit,
	}

	if options.Since > 0 {
		params.Oldest = fmt.Sprintf("%d", options.Since)
	}
	if options.Until > 0 {
		params.Latest = fmt.Sprintf("%d", options.Until)
	}

	history, err := s.client.GetConversationHistory(params)
	if err != nil {
		return nil, fmt.Errorf("failed to get DM history: %w", err)
	}

	messages := make([]MessageInfo, len(history.Messages))
	for i, msg := range history.Messages {
		messages[i] = convertMessage(msg)
	}

	return messages, nil
}

// SendDM sends a direct message to a user
func (s *DMService) SendDM(userID, text string, unfurlLinks, unfurlMedia bool) (string, string, error) {
	// Open/get the DM conversation
	conv, _, _, err := s.client.OpenConversation(&slack.OpenConversationParameters{
		Users: []string{userID},
	})
	if err != nil {
		return "", "", fmt.Errorf("failed to open DM: %w", err)
	}

	// Build message options
	msgOptions := []slack.MsgOption{
		slack.MsgOptionText(text, false),
	}
	if !unfurlLinks {
		msgOptions = append(msgOptions, slack.MsgOptionDisableLinkUnfurl())
	}
	if !unfurlMedia {
		msgOptions = append(msgOptions, slack.MsgOptionDisableMediaUnfurl())
	}

	// Send message
	channel, timestamp, err := s.client.PostMessage(conv.ID, msgOptions...)
	if err != nil {
		return "", "", fmt.Errorf("failed to send DM: %w", err)
	}

	return channel, timestamp, nil
}

// ReplyInDM replies to a message in a DM thread
func (s *DMService) ReplyInDM(userID, threadTS, text string, unfurlLinks, unfurlMedia bool) (string, string, error) {
	// Open/get the DM conversation
	conv, _, _, err := s.client.OpenConversation(&slack.OpenConversationParameters{
		Users: []string{userID},
	})
	if err != nil {
		return "", "", fmt.Errorf("failed to open DM: %w", err)
	}

	// Build message options with thread_ts
	msgOptions := []slack.MsgOption{
		slack.MsgOptionText(text, false),
		slack.MsgOptionTS(threadTS),
	}
	if !unfurlLinks {
		msgOptions = append(msgOptions, slack.MsgOptionDisableLinkUnfurl())
	}
	if !unfurlMedia {
		msgOptions = append(msgOptions, slack.MsgOptionDisableMediaUnfurl())
	}

	// Send message
	channel, timestamp, err := s.client.PostMessage(conv.ID, msgOptions...)
	if err != nil {
		return "", "", fmt.Errorf("failed to reply in DM: %w", err)
	}

	return channel, timestamp, nil
}
