package slack

import (
	"github.com/slack-go/slack"
	"github.com/stretchr/testify/mock"
)

// MockClient is a mock implementation of the Client interface
type MockClient struct {
	mock.Mock
}

func (m *MockClient) GetConversations(params *slack.GetConversationsParameters) ([]slack.Channel, string, error) {
	args := m.Called(params)
	return args.Get(0).([]slack.Channel), args.String(1), args.Error(2)
}

func (m *MockClient) GetConversationInfo(input *slack.GetConversationInfoInput) (*slack.Channel, error) {
	args := m.Called(input)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*slack.Channel), args.Error(1)
}

func (m *MockClient) GetConversationHistory(params *slack.GetConversationHistoryParameters) (*slack.GetConversationHistoryResponse, error) {
	args := m.Called(params)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*slack.GetConversationHistoryResponse), args.Error(1)
}

func (m *MockClient) GetConversationReplies(params *slack.GetConversationRepliesParameters) ([]slack.Message, bool, string, error) {
	args := m.Called(params)
	return args.Get(0).([]slack.Message), args.Bool(1), args.String(2), args.Error(3)
}

func (m *MockClient) GetUsersInConversation(params *slack.GetUsersInConversationParameters) ([]string, string, error) {
	args := m.Called(params)
	return args.Get(0).([]string), args.String(1), args.Error(2)
}

func (m *MockClient) GetUsers(options ...slack.GetUsersOption) ([]slack.User, error) {
	args := m.Called(options)
	return args.Get(0).([]slack.User), args.Error(1)
}

func (m *MockClient) GetUserByEmail(email string) (*slack.User, error) {
	args := m.Called(email)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*slack.User), args.Error(1)
}

func (m *MockClient) GetUserInfo(user string) (*slack.User, error) {
	args := m.Called(user)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*slack.User), args.Error(1)
}

func (m *MockClient) PostMessage(channelID string, options ...slack.MsgOption) (string, string, error) {
	args := m.Called(channelID, options)
	return args.String(0), args.String(1), args.Error(2)
}

func (m *MockClient) UpdateMessage(channelID, timestamp string, options ...slack.MsgOption) (string, string, string, error) {
	args := m.Called(channelID, timestamp, options)
	return args.String(0), args.String(1), args.String(2), args.Error(3)
}

func (m *MockClient) ScheduleMessage(channelID, postAt string, options ...slack.MsgOption) (string, string, error) {
	args := m.Called(channelID, postAt, options)
	return args.String(0), args.String(1), args.Error(2)
}

func (m *MockClient) GetReactions(item slack.ItemRef, params slack.GetReactionsParameters) ([]slack.ItemReaction, error) {
	args := m.Called(item, params)
	return args.Get(0).([]slack.ItemReaction), args.Error(1)
}

func (m *MockClient) AddReaction(name string, item slack.ItemRef) error {
	args := m.Called(name, item)
	return args.Error(0)
}

func (m *MockClient) RemoveReaction(name string, item slack.ItemRef) error {
	args := m.Called(name, item)
	return args.Error(0)
}

func (m *MockClient) CreateConversation(params slack.CreateConversationParams) (*slack.Channel, error) {
	args := m.Called(params)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*slack.Channel), args.Error(1)
}

func (m *MockClient) ArchiveConversation(channelID string) error {
	args := m.Called(channelID)
	return args.Error(0)
}

func (m *MockClient) UnArchiveConversation(channelID string) error {
	args := m.Called(channelID)
	return args.Error(0)
}

func (m *MockClient) RenameConversation(channelID, name string) (*slack.Channel, error) {
	args := m.Called(channelID, name)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*slack.Channel), args.Error(1)
}

func (m *MockClient) SetTopicOfConversation(channelID, topic string) (*slack.Channel, error) {
	args := m.Called(channelID, topic)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*slack.Channel), args.Error(1)
}

func (m *MockClient) SetPurposeOfConversation(channelID, purpose string) (*slack.Channel, error) {
	args := m.Called(channelID, purpose)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*slack.Channel), args.Error(1)
}

func (m *MockClient) InviteUsersToConversation(channelID string, users ...string) (*slack.Channel, error) {
	args := m.Called(channelID, users)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*slack.Channel), args.Error(1)
}

func (m *MockClient) KickUserFromConversation(channelID, user string) error {
	args := m.Called(channelID, user)
	return args.Error(0)
}

func (m *MockClient) OpenConversation(params *slack.OpenConversationParameters) (*slack.Channel, bool, bool, error) {
	args := m.Called(params)
	if args.Get(0) == nil {
		return nil, args.Bool(1), args.Bool(2), args.Error(3)
	}
	return args.Get(0).(*slack.Channel), args.Bool(1), args.Bool(2), args.Error(3)
}

func (m *MockClient) MarkConversation(channelID, timestamp string) error {
	args := m.Called(channelID, timestamp)
	return args.Error(0)
}
