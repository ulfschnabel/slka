package slack

import (
	"errors"
	"testing"

	"github.com/slack-go/slack"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestListDMsSuccess(t *testing.T) {
	mockClient := new(MockClient)
	mockClient.On("GetConversations", mock.MatchedBy(func(params *slack.GetConversationsParameters) bool {
		return params.Types[0] == "im"
	})).Return(
		[]slack.Channel{
			{
				GroupConversation: slack.GroupConversation{
					Conversation: slack.Conversation{
						ID:        "D123",
						IsIM:      true,
						User:      "U456",
						IsPrivate: true,
					},
				},
			},
			{
				GroupConversation: slack.GroupConversation{
					Conversation: slack.Conversation{
						ID:        "D124",
						IsIM:      true,
						User:      "U789",
						IsPrivate: true,
					},
				},
			},
		},
		"",
		nil,
	)

	// Mock user lookups
	mockClient.On("GetUserInfo", "U456").Return(&slack.User{
		ID:   "U456",
		Name: "alice",
	}, nil)
	mockClient.On("GetUserInfo", "U789").Return(&slack.User{
		ID:   "U789",
		Name: "bob",
	}, nil)

	svc := NewDMService(mockClient)
	result, err := svc.List(0)

	assert.NoError(t, err)
	assert.Len(t, result, 2)
	assert.Equal(t, "D123", result[0].ID)
	assert.Len(t, result[0].UserIDs, 1)
	assert.Equal(t, "U456", result[0].UserIDs[0])
	assert.Len(t, result[0].UserNames, 1)
	assert.Equal(t, "alice", result[0].UserNames[0])
	mockClient.AssertExpectations(t)
}

func TestListDMsEmpty(t *testing.T) {
	mockClient := new(MockClient)
	mockClient.On("GetConversations", mock.Anything).Return(
		[]slack.Channel{},
		"",
		nil,
	)

	svc := NewDMService(mockClient)
	result, err := svc.List(0)

	assert.NoError(t, err)
	assert.Len(t, result, 0)
	mockClient.AssertExpectations(t)
}

func TestResolveUserByID(t *testing.T) {
	mockClient := new(MockClient)
	// Not called for ID format
	svc := NewDMService(mockClient)
	userID, err := svc.ResolveUser("U123")

	assert.NoError(t, err)
	assert.Equal(t, "U123", userID)
	mockClient.AssertExpectations(t)
}

func TestResolveUserByEmail(t *testing.T) {
	mockClient := new(MockClient)
	mockClient.On("GetUserByEmail", "alice@example.com").Return(&slack.User{
		ID: "U456",
	}, nil)

	svc := NewDMService(mockClient)
	userID, err := svc.ResolveUser("alice@example.com")

	assert.NoError(t, err)
	assert.Equal(t, "U456", userID)
	mockClient.AssertExpectations(t)
}

func TestResolveUserByName(t *testing.T) {
	mockClient := new(MockClient)
	mockClient.On("GetUserByEmail", "alice").Return(
		(*slack.User)(nil),
		errors.New("user_not_found"),
	)
	mockClient.On("GetUsers", mock.Anything).Return(
		[]slack.User{
			{ID: "U123", Name: "bob"},
			{ID: "U456", Name: "alice"},
			{ID: "U789", Name: "charlie"},
		},
		nil,
	)

	svc := NewDMService(mockClient)
	userID, err := svc.ResolveUser("alice")

	assert.NoError(t, err)
	assert.Equal(t, "U456", userID)
	mockClient.AssertExpectations(t)
}

func TestResolveUserNotFound(t *testing.T) {
	mockClient := new(MockClient)
	mockClient.On("GetUserByEmail", "nonexistent").Return(
		(*slack.User)(nil),
		errors.New("user_not_found"),
	)
	mockClient.On("GetUsers", mock.Anything).Return(
		[]slack.User{
			{ID: "U123", Name: "bob"},
		},
		nil,
	)

	svc := NewDMService(mockClient)
	_, err := svc.ResolveUser("nonexistent")

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "user not found")
	mockClient.AssertExpectations(t)
}

func TestGetDMHistorySuccess(t *testing.T) {
	mockClient := new(MockClient)

	// First resolve user
	mockClient.On("GetUserByEmail", "alice").Return(
		(*slack.User)(nil),
		errors.New("user_not_found"),
	)
	mockClient.On("GetUsers", mock.Anything).Return(
		[]slack.User{
			{ID: "U456", Name: "alice"},
		},
		nil,
	)

	// Open/get DM conversation
	mockClient.On("OpenConversation", mock.MatchedBy(func(params *slack.OpenConversationParameters) bool {
		return params.Users[0] == "U456"
	})).Return(
		&slack.Channel{
			GroupConversation: slack.GroupConversation{
				Conversation: slack.Conversation{
					ID: "D123",
				},
			},
		},
		false, // noOp
		false, // alreadyOpen
		nil,
	)

	// Get history
	mockClient.On("GetConversationHistory", mock.MatchedBy(func(params *slack.GetConversationHistoryParameters) bool {
		return params.ChannelID == "D123"
	})).Return(
		&slack.GetConversationHistoryResponse{
			Messages: []slack.Message{
				{
					Msg: slack.Msg{
						Timestamp: "1706123456.789000",
						User:      "U456",
						Text:      "Hello!",
					},
				},
				{
					Msg: slack.Msg{
						Timestamp: "1706123457.000000",
						User:      "U789",
						Text:      "Hi there!",
					},
				},
			},
		},
		nil,
	)

	svc := NewDMService(mockClient)
	history, err := svc.GetHistory("alice", HistoryOptions{Limit: 10})

	assert.NoError(t, err)
	assert.Len(t, history, 2)
	assert.Equal(t, "Hello!", history[0].Text)
	mockClient.AssertExpectations(t)
}

func TestSendDMSuccess(t *testing.T) {
	mockClient := new(MockClient)

	// Resolve user
	svc := NewDMService(mockClient)

	// Open conversation
	mockClient.On("OpenConversation", mock.MatchedBy(func(params *slack.OpenConversationParameters) bool {
		return len(params.Users) == 1 && params.Users[0] == "U456"
	})).Return(
		&slack.Channel{
			GroupConversation: slack.GroupConversation{
				Conversation: slack.Conversation{
					ID: "D123",
				},
			},
		},
		false,
		false,
		nil,
	)

	// Send message
	mockClient.On("PostMessage", "D123", mock.Anything).Return(
		"D123",
		"1706123456.789000",
		nil,
	)

	channelID, timestamp, err := svc.SendDM([]string{"U456"}, "Hello there!", false, false)

	assert.NoError(t, err)
	assert.Equal(t, "D123", channelID)
	assert.Equal(t, "1706123456.789000", timestamp)
	mockClient.AssertExpectations(t)
}

func TestSendDMOpenConversationFails(t *testing.T) {
	mockClient := new(MockClient)

	mockClient.On("OpenConversation", mock.Anything).Return(
		(*slack.Channel)(nil),
		false,
		false,
		errors.New("user_not_found"),
	)

	svc := NewDMService(mockClient)
	_, _, err := svc.SendDM([]string{"U456"}, "Hello", false, false)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "failed to open DM")
	mockClient.AssertExpectations(t)
}

func TestResolveUsersMultiple(t *testing.T) {
	mockClient := new(MockClient)

	// Mock user lookups
	mockClient.On("GetUserByEmail", "alice").Return(
		(*slack.User)(nil),
		errors.New("not_found"),
	)
	mockClient.On("GetUsers", mock.Anything).Return(
		[]slack.User{
			{ID: "U123", Name: "alice"},
			{ID: "U456", Name: "bob"},
		},
		nil,
	).Once()

	mockClient.On("GetUserByEmail", "bob").Return(
		(*slack.User)(nil),
		errors.New("not_found"),
	)
	mockClient.On("GetUsers", mock.Anything).Return(
		[]slack.User{
			{ID: "U123", Name: "alice"},
			{ID: "U456", Name: "bob"},
		},
		nil,
	).Once()

	svc := NewDMService(mockClient)
	userIDs, err := svc.ResolveUsers("alice,bob")

	assert.NoError(t, err)
	assert.Len(t, userIDs, 2)
	assert.Contains(t, userIDs, "U123")
	assert.Contains(t, userIDs, "U456")
	mockClient.AssertExpectations(t)
}

func TestSendGroupDM(t *testing.T) {
	mockClient := new(MockClient)

	// Open group conversation
	mockClient.On("OpenConversation", mock.MatchedBy(func(params *slack.OpenConversationParameters) bool {
		return len(params.Users) == 3
	})).Return(
		&slack.Channel{
			GroupConversation: slack.GroupConversation{
				Conversation: slack.Conversation{
					ID:     "G123",
					IsMpIM: true,
				},
			},
		},
		false,
		false,
		nil,
	)

	// Send message
	mockClient.On("PostMessage", "G123", mock.Anything).Return(
		"G123",
		"1706123456.789000",
		nil,
	)

	svc := NewDMService(mockClient)
	channelID, timestamp, err := svc.SendDM([]string{"U123", "U456", "U789"}, "Hello team!", false, false)

	assert.NoError(t, err)
	assert.Equal(t, "G123", channelID)
	assert.Equal(t, "1706123456.789000", timestamp)
	mockClient.AssertExpectations(t)
}

func TestReplyInDMSuccess(t *testing.T) {
	mockClient := new(MockClient)

	// Open conversation
	mockClient.On("OpenConversation", mock.Anything).Return(
		&slack.Channel{
			GroupConversation: slack.GroupConversation{
				Conversation: slack.Conversation{
					ID: "D123",
				},
			},
		},
		false,
		false,
		nil,
	)

	// Post message with thread_ts
	mockClient.On("PostMessage", "D123", mock.Anything).Return(
		"D123",
		"1706123457.000000",
		nil,
	)

	svc := NewDMService(mockClient)
	channelID, timestamp, err := svc.ReplyInDM([]string{"U456"}, "1706123456.789000", "Replying!", false, false)

	assert.NoError(t, err)
	assert.Equal(t, "D123", channelID)
	assert.Equal(t, "1706123457.000000", timestamp)
	mockClient.AssertExpectations(t)
}
