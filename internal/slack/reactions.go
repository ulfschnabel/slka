package slack

import (
	"fmt"
	"strings"

	"github.com/slack-go/slack"
)

// AcknowledgmentInfo represents acknowledgment status of a message
type AcknowledgmentInfo struct {
	IsAcknowledged bool     `json:"is_acknowledged"`
	ReactedUsers   []string `json:"reacted_users,omitempty"`
	ReactionCount  int      `json:"reaction_count"`
	ReplyCount     int      `json:"reply_count"`
	HasReplies     bool     `json:"has_replies"`
	HasReactions   bool     `json:"has_reactions"`
	MessageAuthor  string   `json:"message_author"`
}

// ReactionListInfo represents detailed reaction information
type ReactionListInfo struct {
	Timestamp          string         `json:"timestamp"`
	Channel            string         `json:"channel"`
	MessageText        string         `json:"message_text,omitempty"`
	Reactions          []ReactionInfo `json:"reactions"`
	TotalReactionCount int            `json:"total_reaction_count"`
}

// ReactionService provides reaction operations
type ReactionService struct {
	client Client
}

// NewReactionService creates a new ReactionService
func NewReactionService(client Client) *ReactionService {
	return &ReactionService{client: client}
}

// AddReaction adds a reaction to a message
func (s *ReactionService) AddReaction(channelID, timestamp, emoji string) error {
	// Strip colons from emoji if present
	emoji = strings.Trim(emoji, ":")

	item := slack.ItemRef{
		Channel:   channelID,
		Timestamp: timestamp,
	}

	if err := s.client.AddReaction(emoji, item); err != nil {
		return fmt.Errorf("failed to add reaction: %w", err)
	}

	return nil
}

// RemoveReaction removes a reaction from a message
func (s *ReactionService) RemoveReaction(channelID, timestamp, emoji string) error {
	// Strip colons from emoji if present
	emoji = strings.Trim(emoji, ":")

	item := slack.ItemRef{
		Channel:   channelID,
		Timestamp: timestamp,
	}

	if err := s.client.RemoveReaction(emoji, item); err != nil {
		return fmt.Errorf("failed to remove reaction: %w", err)
	}

	return nil
}

// ListReactions gets all reactions on a message
func (s *ReactionService) ListReactions(channelID, timestamp string) (*ReactionListInfo, error) {
	// Get reactions
	item := slack.ItemRef{
		Channel:   channelID,
		Timestamp: timestamp,
	}

	reactions, err := s.client.GetReactions(item, slack.GetReactionsParameters{})
	if err != nil {
		return nil, fmt.Errorf("failed to get reactions: %w", err)
	}

	// Get message text for context
	params := &slack.GetConversationHistoryParameters{
		ChannelID: channelID,
		Latest:    timestamp,
		Limit:     1,
		Inclusive: true,
	}

	history, err := s.client.GetConversationHistory(params)
	if err != nil {
		return nil, fmt.Errorf("failed to get message: %w", err)
	}

	messageText := ""
	if len(history.Messages) > 0 {
		messageText = history.Messages[0].Text
	}

	// Convert to ReactionInfo and calculate total
	reactionInfos := make([]ReactionInfo, len(reactions))
	totalCount := 0
	for i, r := range reactions {
		reactionInfos[i] = ReactionInfo{
			Name:  r.Name,
			Count: r.Count,
			Users: r.Users,
		}
		totalCount += r.Count
	}

	return &ReactionListInfo{
		Timestamp:          timestamp,
		Channel:            channelID,
		MessageText:        messageText,
		Reactions:          reactionInfos,
		TotalReactionCount: totalCount,
	}, nil
}

// CheckAcknowledgment checks if a message has been acknowledged by others
// Acknowledgment = any reaction from non-author OR any reply from non-author
func (s *ReactionService) CheckAcknowledgment(channelID, timestamp string) (*AcknowledgmentInfo, error) {
	// Get message with reactions
	params := &slack.GetConversationHistoryParameters{
		ChannelID: channelID,
		Latest:    timestamp,
		Limit:     1,
		Inclusive: true,
	}

	history, err := s.client.GetConversationHistory(params)
	if err != nil {
		return nil, fmt.Errorf("failed to get message: %w", err)
	}

	if len(history.Messages) == 0 {
		return nil, fmt.Errorf("message not found")
	}

	msg := history.Messages[0]
	author := msg.User

	// Check reactions from others
	reactedUsers := []string{}
	reactionCount := 0
	hasReactions := len(msg.Reactions) > 0

	for _, reaction := range msg.Reactions {
		for _, userID := range reaction.Users {
			if userID != author {
				reactedUsers = append(reactedUsers, userID)
				reactionCount++
			}
		}
	}

	// Check replies from others
	replyCount := 0
	hasReplies := msg.ReplyCount > 0

	if msg.ReplyCount > 0 {
		replyParams := &slack.GetConversationRepliesParameters{
			ChannelID: channelID,
			Timestamp: timestamp,
		}

		replies, _, _, err := s.client.GetConversationReplies(replyParams)
		if err != nil {
			return nil, fmt.Errorf("failed to get replies: %w", err)
		}

		// Count replies from others (skip first message which is the parent)
		for i, reply := range replies {
			if i == 0 {
				continue // Skip parent message
			}
			if reply.User != author {
				replyCount++
			}
		}
	}

	isAcknowledged := reactionCount > 0 || replyCount > 0

	return &AcknowledgmentInfo{
		IsAcknowledged: isAcknowledged,
		ReactedUsers:   reactedUsers,
		ReactionCount:  reactionCount,
		ReplyCount:     replyCount,
		HasReplies:     hasReplies,
		HasReactions:   hasReactions,
		MessageAuthor:  author,
	}, nil
}
