package slack

import (
	"fmt"
	"sort"
	"strings"

	"github.com/slack-go/slack"
)

// DMInfo represents a direct message conversation
type DMInfo struct {
	ID            string   `json:"id"`
	Type          string   `json:"type"` // "im" or "mpim"
	UserIDs       []string `json:"user_ids"`
	UserNames     []string `json:"user_names,omitempty"`
	LastMessageTS string   `json:"last_message_ts,omitempty"`
}

// DMService provides DM operations
type DMService struct {
	client Client
}

// NewDMService creates a new DMService
func NewDMService(client Client) *DMService {
	return &DMService{client: client}
}

// List returns all DM conversations (both 1-on-1 and group)
func (s *DMService) List(limit int) ([]DMInfo, error) {
	// Get both im and mpim conversations
	params := &slack.GetConversationsParameters{
		Types: []string{"im", "mpim"},
		Limit: limit,
	}

	channels, _, err := s.client.GetConversations(params)
	if err != nil {
		return nil, fmt.Errorf("failed to list DMs: %w", err)
	}

	// Collect all user IDs we need to resolve
	var userIDsToResolve []string
	for _, channel := range channels {
		if channel.IsIM && channel.User != "" {
			userIDsToResolve = append(userIDsToResolve, channel.User)
		} else if channel.IsMpIM {
			userIDsToResolve = append(userIDsToResolve, channel.Members...)
		}
	}

	// Batch resolve user names (N calls, but bounded by limit)
	userMap := ResolveUserNames(s.client, userIDsToResolve)

	dms := make([]DMInfo, 0, len(channels))
	for _, channel := range channels {
		dm := DMInfo{
			ID:   channel.ID,
			Type: "im",
		}

		if channel.Latest != nil {
			dm.LastMessageTS = channel.Latest.Timestamp
		}

		// For regular DMs (im type)
		if channel.IsIM {
			dm.UserIDs = []string{channel.User}
			if info, ok := userMap[channel.User]; ok {
				dm.UserNames = []string{info.Name}
			}
		} else if channel.IsMpIM {
			// For group DMs (mpim type)
			dm.Type = "mpim"
			if len(channel.Members) > 0 {
				dm.UserIDs = channel.Members
				for _, userID := range channel.Members {
					if info, ok := userMap[userID]; ok {
						dm.UserNames = append(dm.UserNames, info.Name)
					}
				}
			}
		}

		dms = append(dms, dm)
	}

	// Sort by last message timestamp descending (most recent first)
	sort.Slice(dms, func(i, j int) bool {
		return dms[i].LastMessageTS > dms[j].LastMessageTS
	})

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

// ResolveUsers converts comma-separated user identifiers to user IDs
func (s *DMService) ResolveUsers(usersArg string) ([]string, error) {
	// Split by comma and trim spaces
	userArgs := strings.Split(usersArg, ",")
	userIDs := make([]string, 0, len(userArgs))

	for _, userArg := range userArgs {
		userArg = strings.TrimSpace(userArg)
		if userArg == "" {
			continue
		}

		userID, err := s.ResolveUser(userArg)
		if err != nil {
			return nil, err
		}
		userIDs = append(userIDs, userID)
	}

	if len(userIDs) == 0 {
		return nil, fmt.Errorf("no valid users specified")
	}

	return userIDs, nil
}

// GetHistory gets the message history of a DM conversation
func (s *DMService) GetHistory(usersArg string, options HistoryOptions) ([]MessageInfo, error) {
	// Resolve users to IDs
	userIDs, err := s.ResolveUsers(usersArg)
	if err != nil {
		return nil, err
	}

	// Open/get the DM conversation (works for both im and mpim)
	conv, _, _, err := s.client.OpenConversation(&slack.OpenConversationParameters{
		Users: userIDs,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to open DM: %w", err)
	}

	// Get conversation history
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

// SendDM sends a direct message to one or more users
func (s *DMService) SendDM(userIDs []string, text string, unfurlLinks, unfurlMedia bool) (string, string, error) {
	// Open/get the DM conversation (creates group DM if multiple users)
	conv, _, _, err := s.client.OpenConversation(&slack.OpenConversationParameters{
		Users: userIDs,
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
func (s *DMService) ReplyInDM(userIDs []string, threadTS, text string, unfurlLinks, unfurlMedia bool) (string, string, error) {
	// Open/get the DM conversation
	conv, _, _, err := s.client.OpenConversation(&slack.OpenConversationParameters{
		Users: userIDs,
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

// FindExistingConversation finds an existing DM/group DM with the specified users
// Returns conversation ID if found, empty string if not found
func (s *DMService) FindExistingConversation(userIDs []string) (string, error) {
	// Sort user IDs for comparison
	sortedUserIDs := make([]string, len(userIDs))
	copy(sortedUserIDs, userIDs)
	sort.Strings(sortedUserIDs)

	// Get all DM conversations
	convType := "im"
	if len(userIDs) > 1 {
		convType = "mpim"
	}

	params := &slack.GetConversationsParameters{
		Types: []string{convType},
	}

	channels, _, err := s.client.GetConversations(params)
	if err != nil {
		return "", fmt.Errorf("failed to list conversations: %w", err)
	}

	// For each conversation, check if user list matches
	for _, channel := range channels {
		var convUsers []string
		if channel.IsIM {
			convUsers = []string{channel.User}
		} else if channel.IsMpIM && len(channel.Members) > 0 {
			convUsers = channel.Members
		}

		// Sort and compare
		if len(convUsers) == len(sortedUserIDs) {
			sort.Strings(convUsers)
			match := true
			for i, userID := range sortedUserIDs {
				if convUsers[i] != userID {
					match = false
					break
				}
			}
			if match {
				return channel.ID, nil
			}
		}
	}

	return "", nil // Not found
}
