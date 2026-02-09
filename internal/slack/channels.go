package slack

import (
	"fmt"
	"strings"

	"github.com/slack-go/slack"
)

// ChannelService handles channel-related operations
type ChannelService struct {
	client Client
}

// NewChannelService creates a new channel service
func NewChannelService(client Client) *ChannelService {
	return &ChannelService{client: client}
}

// ChannelInfo represents channel information
type ChannelInfo struct {
	ID           string `json:"id"`
	Name         string `json:"name"`
	IsPrivate    bool   `json:"is_private"`
	IsArchived   bool   `json:"is_archived"`
	Topic        string `json:"topic,omitempty"`
	Purpose      string `json:"purpose,omitempty"`
	MemberCount  int    `json:"member_count,omitempty"`
	Created      int64  `json:"created"`
	Creator      string `json:"creator,omitempty"`
	LastMessageTS string `json:"last_message_ts,omitempty"`
}

// MessageInfo represents a message
type MessageInfo struct {
	Timestamp   string                   `json:"ts"`
	User        string                   `json:"user"`
	UserName    string                   `json:"user_name,omitempty"`
	Text        string                   `json:"text"`
	ThreadTS    string                   `json:"thread_ts,omitempty"`
	ReplyCount  int                      `json:"reply_count"`
	Reactions   []ReactionInfo           `json:"reactions,omitempty"`
	Links       []map[string]interface{} `json:"links,omitempty"`
}

// ReactionInfo represents a reaction
type ReactionInfo struct {
	Name  string   `json:"name"`
	Count int      `json:"count"`
	Users []string `json:"users"`
}

// UserInfo represents a user
type UserInfo struct {
	ID       string `json:"id"`
	Name     string `json:"name"`
	RealName string `json:"real_name"`
	Email    string `json:"email,omitempty"`
	IsBot    bool   `json:"is_bot"`
}

// ListChannelsOptions contains options for listing channels
type ListChannelsOptions struct {
	IncludeArchived bool
	Type            string // "public", "private", "all"
	Limit           int
}

// HistoryOptions contains options for fetching message history
type HistoryOptions struct {
	Since          int64
	Until          int64
	Limit          int
	IncludeThreads bool
}

// List returns all channels
func (s *ChannelService) List(opts ListChannelsOptions) ([]ChannelInfo, error) {
	params := &slack.GetConversationsParameters{
		ExcludeArchived: !opts.IncludeArchived,
		Limit:           opts.Limit,
	}

	if opts.Limit == 0 {
		params.Limit = 1000
	}

	// Set types based on filter
	switch opts.Type {
	case "public":
		params.Types = []string{"public_channel"}
	case "private":
		params.Types = []string{"private_channel"}
	default:
		params.Types = []string{"public_channel", "private_channel"}
	}

	channels, _, err := s.client.GetConversations(params)
	if err != nil {
		return nil, fmt.Errorf("failed to get conversations: %w", err)
	}

	result := make([]ChannelInfo, len(channels))
	for i, ch := range channels {
		result[i] = convertChannel(ch)
	}

	return result, nil
}

// GetInfo returns detailed information about a channel
func (s *ChannelService) GetInfo(channelID string) (*ChannelInfo, error) {
	ch, err := s.client.GetConversationInfo(&slack.GetConversationInfoInput{
		ChannelID: channelID,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to get channel info: %w", err)
	}

	info := convertChannel(*ch)
	return &info, nil
}

// GetHistory returns message history for a channel
func (s *ChannelService) GetHistory(channelID string, opts HistoryOptions) ([]MessageInfo, error) {
	params := &slack.GetConversationHistoryParameters{
		ChannelID: channelID,
		Limit:     opts.Limit,
	}

	if opts.Limit == 0 {
		params.Limit = 100
	}

	if opts.Since > 0 {
		params.Oldest = fmt.Sprintf("%d", opts.Since)
	}

	if opts.Until > 0 {
		params.Latest = fmt.Sprintf("%d", opts.Until)
	}

	resp, err := s.client.GetConversationHistory(params)
	if err != nil {
		return nil, fmt.Errorf("failed to get conversation history: %w", err)
	}

	messages := make([]MessageInfo, len(resp.Messages))
	for i, msg := range resp.Messages {
		messages[i] = convertMessage(msg)
	}

	return messages, nil
}

// GetMembers returns members of a channel
func (s *ChannelService) GetMembers(channelID string, limit int) ([]UserInfo, error) {
	if limit == 0 {
		limit = 1000
	}

	params := &slack.GetUsersInConversationParameters{
		ChannelID: channelID,
		Limit:     limit,
	}

	userIDs, _, err := s.client.GetUsersInConversation(params)
	if err != nil {
		return nil, fmt.Errorf("failed to get users in conversation: %w", err)
	}

	users := make([]UserInfo, 0, len(userIDs))
	for _, userID := range userIDs {
		user, err := s.client.GetUserInfo(userID)
		if err != nil {
			// Skip users we can't fetch info for
			continue
		}
		users = append(users, convertUser(*user))
	}

	return users, nil
}

// ResolveChannel resolves a channel name or ID to a channel ID
func (s *ChannelService) ResolveChannel(channel string) (string, error) {
	// If it starts with C or G, it's probably an ID
	if strings.HasPrefix(channel, "C") || strings.HasPrefix(channel, "G") {
		return channel, nil
	}

	// Remove # prefix if present
	name := strings.TrimPrefix(channel, "#")

	// List all channels and find matching name
	channels, err := s.List(ListChannelsOptions{
		IncludeArchived: true,
		Type:            "all",
	})
	if err != nil {
		return "", err
	}

	for _, ch := range channels {
		if ch.Name == name {
			return ch.ID, nil
		}
	}

	return "", fmt.Errorf("channel not found: %s", channel)
}

// Helper functions

func convertChannel(ch slack.Channel) ChannelInfo {
	return ChannelInfo{
		ID:          ch.ID,
		Name:        ch.Name,
		IsPrivate:   ch.IsPrivate,
		IsArchived:  ch.IsArchived,
		Topic:       ch.Topic.Value,
		Purpose:     ch.Purpose.Value,
		MemberCount: ch.NumMembers,
		Created:     int64(ch.Created),
		Creator:     ch.Creator,
	}
}

// MarkAsRead marks a channel as read up to the specified timestamp
func (s *ChannelService) MarkAsRead(channelID string, timestamp string) error {
	err := s.client.MarkConversation(channelID, timestamp)
	if err != nil {
		return fmt.Errorf("failed to mark as read: %w", err)
	}
	return nil
}

func convertMessage(msg slack.Message) MessageInfo {
	reactions := make([]ReactionInfo, len(msg.Reactions))
	for i, r := range msg.Reactions {
		reactions[i] = ReactionInfo{
			Name:  r.Name,
			Count: r.Count,
			Users: r.Users,
		}
	}

	return MessageInfo{
		Timestamp:  msg.Timestamp,
		User:       msg.User,
		Text:       msg.Text,
		ThreadTS:   msg.ThreadTimestamp,
		ReplyCount: msg.ReplyCount,
		Reactions:  reactions,
	}
}

func convertUser(user slack.User) UserInfo {
	return UserInfo{
		ID:       user.ID,
		Name:     user.Name,
		RealName: user.RealName,
		Email:    user.Profile.Email,
		IsBot:    user.IsBot,
	}
}
