package slack

import (
	"testing"
	"time"

	"github.com/slack-go/slack"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestListChannelsSuccess(t *testing.T) {
	mockClient := new(MockClient)
	mockClient.On("GetConversations", mock.Anything).Return(
		[]slack.Channel{
			{
				GroupConversation: slack.GroupConversation{
					Conversation: slack.Conversation{
						ID:        "C123",
						IsPrivate: false,
					},
					Name: "general",
				},
			},
		},
		"",  // cursor
		nil, // error
	)

	svc := NewChannelService(mockClient)
	result, err := svc.List(ListChannelsOptions{})

	assert.NoError(t, err)
	assert.Len(t, result, 1)
	assert.Equal(t, "general", result[0].Name)
	assert.Equal(t, "C123", result[0].ID)
	mockClient.AssertExpectations(t)
}

func TestListChannelsExcludesArchivedByDefault(t *testing.T) {
	mockClient := new(MockClient)
	mockClient.On("GetConversations", mock.MatchedBy(func(params *slack.GetConversationsParameters) bool {
		return params.ExcludeArchived == true
	})).Return([]slack.Channel{}, "", nil)

	svc := NewChannelService(mockClient)
	_, err := svc.List(ListChannelsOptions{IncludeArchived: false})

	assert.NoError(t, err)
	mockClient.AssertExpectations(t)
}

func TestListChannelsIncludesArchived(t *testing.T) {
	mockClient := new(MockClient)
	mockClient.On("GetConversations", mock.MatchedBy(func(params *slack.GetConversationsParameters) bool {
		return params.ExcludeArchived == false
	})).Return([]slack.Channel{}, "", nil)

	svc := NewChannelService(mockClient)
	_, err := svc.List(ListChannelsOptions{IncludeArchived: true})

	assert.NoError(t, err)
	mockClient.AssertExpectations(t)
}

func TestGetChannelInfo(t *testing.T) {
	mockClient := new(MockClient)
	mockClient.On("GetConversationInfo", mock.Anything).Return(
		&slack.Channel{
			GroupConversation: slack.GroupConversation{
				Conversation: slack.Conversation{
					ID: "C123",
				},
				Name:    "general",
				Topic:   slack.Topic{Value: "General discussion"},
				Purpose: slack.Purpose{Value: "Company-wide"},
			},
		},
		nil,
	)

	svc := NewChannelService(mockClient)
	result, err := svc.GetInfo("C123")

	assert.NoError(t, err)
	assert.Equal(t, "C123", result.ID)
	assert.Equal(t, "general", result.Name)
	assert.Equal(t, "General discussion", result.Topic)
	mockClient.AssertExpectations(t)
}

func TestGetChannelHistory(t *testing.T) {
	ts := time.Now().Unix()
	mockClient := new(MockClient)
	mockClient.On("GetConversationHistory", mock.Anything).Return(
		&slack.GetConversationHistoryResponse{
			Messages: []slack.Message{
				{
					Msg: slack.Msg{
						Timestamp: "1706123456.789000",
						Text:      "Hello world",
						User:      "U123",
					},
				},
			},
		},
		nil,
	)

	svc := NewChannelService(mockClient)
	result, err := svc.GetHistory("C123", HistoryOptions{
		Since: ts,
		Limit: 100,
	})

	assert.NoError(t, err)
	assert.Len(t, result, 1)
	assert.Equal(t, "Hello world", result[0].Text)
	mockClient.AssertExpectations(t)
}

func TestGetChannelMembers(t *testing.T) {
	mockClient := new(MockClient)
	mockClient.On("GetUsersInConversation", mock.Anything).Return(
		[]string{"U123", "U456"},
		"", // cursor
		nil,
	)

	// Mock user info calls
	mockClient.On("GetUserInfo", "U123").Return(
		&slack.User{
			ID:       "U123",
			Name:     "johndoe",
			RealName: "John Doe",
		},
		nil,
	)
	mockClient.On("GetUserInfo", "U456").Return(
		&slack.User{
			ID:       "U456",
			Name:     "janedoe",
			RealName: "Jane Doe",
		},
		nil,
	)

	svc := NewChannelService(mockClient)
	result, err := svc.GetMembers("C123", 100)

	assert.NoError(t, err)
	assert.Len(t, result, 2)
	assert.Equal(t, "johndoe", result[0].Name)
	assert.Equal(t, "janedoe", result[1].Name)
	mockClient.AssertExpectations(t)
}

func TestListChannelsSortedByLastMessage(t *testing.T) {
	mockClient := new(MockClient)
	mockClient.On("GetConversations", mock.Anything).Return(
		[]slack.Channel{
			{
				GroupConversation: slack.GroupConversation{
					Conversation: slack.Conversation{
						ID: "C_OLD",
					},
					Name: "old-channel",
				},
				// No Latest message
			},
			{
				GroupConversation: slack.GroupConversation{
					Conversation: slack.Conversation{
						ID: "C_MID",
						Latest: &slack.Message{
							Msg: slack.Msg{Timestamp: "1700000000.000000"},
						},
					},
					Name: "mid-channel",
				},
			},
			{
				GroupConversation: slack.GroupConversation{
					Conversation: slack.Conversation{
						ID: "C_NEW",
						Latest: &slack.Message{
							Msg: slack.Msg{Timestamp: "1706000000.000000"},
						},
					},
					Name: "new-channel",
				},
			},
		},
		"",
		nil,
	)

	svc := NewChannelService(mockClient)
	result, err := svc.List(ListChannelsOptions{})

	assert.NoError(t, err)
	assert.Len(t, result, 3)
	// Most recent first
	assert.Equal(t, "new-channel", result[0].Name)
	assert.Equal(t, "1706000000.000000", result[0].LastMessageTS)
	assert.Equal(t, "mid-channel", result[1].Name)
	assert.Equal(t, "1700000000.000000", result[1].LastMessageTS)
	// No messages sorts last
	assert.Equal(t, "old-channel", result[2].Name)
	assert.Equal(t, "", result[2].LastMessageTS)
	mockClient.AssertExpectations(t)
}

func TestResolveChannelByName(t *testing.T) {
	mockClient := new(MockClient)
	mockClient.On("GetConversations", mock.Anything).Return(
		[]slack.Channel{
			{
				GroupConversation: slack.GroupConversation{
					Conversation: slack.Conversation{
						ID: "C123",
					},
					Name: "general",
				},
			},
		},
		"",
		nil,
	)

	svc := NewChannelService(mockClient)
	id, err := svc.ResolveChannel("general")

	assert.NoError(t, err)
	assert.Equal(t, "C123", id)
	mockClient.AssertExpectations(t)
}

func TestResolveChannelByNameWithHash(t *testing.T) {
	mockClient := new(MockClient)
	mockClient.On("GetConversations", mock.Anything).Return(
		[]slack.Channel{
			{
				GroupConversation: slack.GroupConversation{
					Conversation: slack.Conversation{
						ID: "C123",
					},
					Name: "general",
				},
			},
		},
		"",
		nil,
	)

	svc := NewChannelService(mockClient)
	id, err := svc.ResolveChannel("#general")

	assert.NoError(t, err)
	assert.Equal(t, "C123", id)
	mockClient.AssertExpectations(t)
}

func TestResolveChannelByID(t *testing.T) {
	mockClient := new(MockClient)

	svc := NewChannelService(mockClient)
	id, err := svc.ResolveChannel("C123456")

	assert.NoError(t, err)
	assert.Equal(t, "C123456", id)
	// No API call should be made for IDs
}

func TestResolveChannelNotFound(t *testing.T) {
	mockClient := new(MockClient)
	mockClient.On("GetConversations", mock.Anything).Return(
		[]slack.Channel{},
		"",
		nil,
	)

	svc := NewChannelService(mockClient)
	_, err := svc.ResolveChannel("nonexistent")

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "not found")
	mockClient.AssertExpectations(t)
}
