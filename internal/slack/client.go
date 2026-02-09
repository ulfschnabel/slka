package slack

import (
	"os"

	"github.com/slack-go/slack"
)

// Client defines the interface for Slack API operations
// This allows us to mock the client in tests
type Client interface {
	// Channels
	GetConversations(params *slack.GetConversationsParameters) ([]slack.Channel, string, error)
	GetConversationInfo(input *slack.GetConversationInfoInput) (*slack.Channel, error)
	GetConversationHistory(params *slack.GetConversationHistoryParameters) (*slack.GetConversationHistoryResponse, error)
	GetConversationReplies(params *slack.GetConversationRepliesParameters) ([]slack.Message, bool, string, error)
	GetUsersInConversation(params *slack.GetUsersInConversationParameters) ([]string, string, error)

	// Users
	GetUsers(options ...slack.GetUsersOption) ([]slack.User, error)
	GetUserByEmail(email string) (*slack.User, error)
	GetUserInfo(user string) (*slack.User, error)

	// Messages
	PostMessage(channelID string, options ...slack.MsgOption) (string, string, error)
	UpdateMessage(channelID, timestamp string, options ...slack.MsgOption) (string, string, string, error)
	ScheduleMessage(channelID, postAt string, options ...slack.MsgOption) (string, string, error)

	// Reactions
	GetReactions(item slack.ItemRef, params slack.GetReactionsParameters) ([]slack.ItemReaction, error)
	AddReaction(name string, item slack.ItemRef) error
	RemoveReaction(name string, item slack.ItemRef) error

	// Channel management
	CreateConversation(params slack.CreateConversationParams) (*slack.Channel, error)
	ArchiveConversation(channelID string) error
	UnArchiveConversation(channelID string) error
	RenameConversation(channelID, name string) (*slack.Channel, error)
	SetTopicOfConversation(channelID, topic string) (*slack.Channel, error)
	SetPurposeOfConversation(channelID, purpose string) (*slack.Channel, error)
	InviteUsersToConversation(channelID string, users ...string) (*slack.Channel, error)
	KickUserFromConversation(channelID, user string) error

	// DMs
	OpenConversation(params *slack.OpenConversationParameters) (*slack.Channel, bool, bool, error)
}

// RealClient wraps the actual slack.Client
type RealClient struct {
	*slack.Client
}

// NewClient creates a new Slack client
func NewClient(token string) Client {
	// Check if custom API URL is set (for testing)
	options := []slack.Option{}
	if apiURL := os.Getenv("SLACK_API_URL"); apiURL != "" {
		options = append(options, slack.OptionAPIURL(apiURL))
	}

	return &RealClient{
		Client: slack.New(token, options...),
	}
}

// Ensure RealClient implements Client interface
var _ Client = (*RealClient)(nil)
