package slack

import (
	"sort"

	"github.com/slack-go/slack"
)

// UnreadInfo contains information about a conversation with unread messages
type UnreadInfo struct {
	ID                 string   `json:"id"`
	Name               string   `json:"name,omitempty"`
	Type               string   `json:"type"` // "channel", "im", "mpim"
	IsChannel          bool     `json:"is_channel"`
	IsPrivate          bool     `json:"is_private"`
	IsIM               bool     `json:"is_im"`
	IsMpIM             bool     `json:"is_mpim"`
	UnreadCount        int      `json:"unread_count"`
	UnreadCountDisplay int      `json:"unread_count_display"`
	LastRead           string   `json:"last_read,omitempty"`
	UserID             string   `json:"user_id,omitempty"`  // For 1-on-1 DMs
	UserName           string   `json:"user_name,omitempty"` // For 1-on-1 DMs
	UserIDs            []string `json:"user_ids,omitempty"`  // For group DMs
}

// UnreadOptions configures how unread items are filtered and ordered
type UnreadOptions struct {
	ChannelsOnly   bool   // Only return channels (not DMs)
	DMsOnly        bool   // Only return DMs (1-on-1 and groups)
	MinUnreadCount int    // Minimum number of unread messages (0 = any)
	OrderBy        string // Ordering: "count" (most unread first), "oldest" (oldest unread first), default: "count"
}

// UnreadService handles retrieving unread conversations
type UnreadService struct {
	client Client
}

// NewUnreadService creates a new UnreadService
func NewUnreadService(client Client) *UnreadService {
	return &UnreadService{client: client}
}

// ListUnread returns all conversations with unread messages, ordered by unread count (descending)
func (s *UnreadService) ListUnread(opts UnreadOptions) ([]UnreadInfo, error) {
	// Get all conversations
	params := &slack.GetConversationsParameters{
		Types:           []string{"public_channel", "private_channel", "im", "mpim"},
		Limit:           1000,
		ExcludeArchived: true,
	}

	conversations, _, err := s.client.GetConversations(params)
	if err != nil {
		return nil, err
	}

	var results []UnreadInfo

	for _, conv := range conversations {
		// GetConversations doesn't always return unread counts, so we need to fetch
		// detailed info for each conversation to get accurate unread counts
		convInfo, err := s.client.GetConversationInfo(&slack.GetConversationInfoInput{
			ChannelID:         conv.ID,
			IncludeNumMembers: false,
		})
		if err != nil {
			// If we can't get info for this conversation, skip it
			continue
		}

		// Skip if no unread messages
		if convInfo.UnreadCount == 0 && convInfo.UnreadCountDisplay == 0 {
			continue
		}

		// Skip if below minimum threshold
		unreadCount := convInfo.UnreadCount
		if unreadCount == 0 {
			unreadCount = convInfo.UnreadCountDisplay
		}
		if opts.MinUnreadCount > 0 && unreadCount < opts.MinUnreadCount {
			continue
		}

		// Apply type filters
		isChannel := convInfo.IsChannel
		isDM := convInfo.IsIM || convInfo.IsMpIM

		if opts.ChannelsOnly && !isChannel {
			continue
		}
		if opts.DMsOnly && !isDM {
			continue
		}

		// Determine type string
		convType := "channel"
		if convInfo.IsIM {
			convType = "im"
		} else if convInfo.IsMpIM {
			convType = "mpim"
		}

		info := UnreadInfo{
			ID:                 convInfo.ID,
			Name:               convInfo.Name,
			Type:               convType,
			IsChannel:          convInfo.IsChannel,
			IsPrivate:          convInfo.IsPrivate,
			IsIM:               convInfo.IsIM,
			IsMpIM:             convInfo.IsMpIM,
			UnreadCount:        convInfo.UnreadCount,
			UnreadCountDisplay: convInfo.UnreadCountDisplay,
			LastRead:           convInfo.LastRead,
		}

		// For 1-on-1 DMs, get user information
		if convInfo.IsIM && convInfo.User != "" {
			user, err := s.client.GetUserInfo(convInfo.User)
			if err == nil {
				info.UserID = user.ID
				info.UserName = user.Name
			}
		}

		// For group DMs, could potentially list users (optional for now)
		// This would require additional API calls

		results = append(results, info)
	}

	// Sort according to OrderBy option
	if opts.OrderBy == "oldest" {
		// Sort by last_read timestamp (ascending), so oldest unread items come first
		sort.Slice(results, func(i, j int) bool {
			return results[i].LastRead < results[j].LastRead
		})
	} else {
		// Default: sort by unread count (descending), so most urgent items come first
		sort.Slice(results, func(i, j int) bool {
			return results[i].UnreadCount > results[j].UnreadCount
		})
	}

	return results, nil
}
