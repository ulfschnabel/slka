package slack

import (
	"testing"

	"github.com/slack-go/slack"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestListUnread(t *testing.T) {
	mockClient := new(MockClient)

	// Create channels manually to match slack.Channel structure
	channel1 := slack.Channel{}
	channel1.ID = "C1"
	channel1.Name = "general"
	channel1.IsChannel = true
	channel1.UnreadCount = 5
	channel1.UnreadCountDisplay = 5
	channel1.LastRead = "1706123450.000000"

	channel2 := slack.Channel{}
	channel2.ID = "C2"
	channel2.Name = "random"
	channel2.IsChannel = true
	channel2.UnreadCount = 0 // No unread messages
	channel2.LastRead = "1706123460.000000"

	dm1 := slack.Channel{}
	dm1.ID = "D1"
	dm1.IsIM = true
	dm1.UnreadCount = 3
	dm1.UnreadCountDisplay = 3
	dm1.LastRead = "1706123470.000000"
	dm1.User = "U123"

	groupDM := slack.Channel{}
	groupDM.ID = "G1"
	groupDM.Name = "mpdm-user1--user2--user3-1"
	groupDM.IsMpIM = true
	groupDM.UnreadCount = 10
	groupDM.UnreadCountDisplay = 10
	groupDM.LastRead = "1706123480.000000"

	channels := []slack.Channel{channel1, channel2, dm1, groupDM}

	mockClient.On("GetConversations", mock.Anything).Return(channels, "", nil)

	// Mock user lookup for DM
	mockClient.On("GetUserInfo", "U123").Return(&slack.User{
		ID:   "U123",
		Name: "alice",
	}, nil)

	svc := NewUnreadService(mockClient)
	results, err := svc.ListUnread(UnreadOptions{})

	assert.NoError(t, err)
	assert.Len(t, results, 3, "Should return only conversations with unread messages")

	// Check that results are ordered by unread count (descending)
	assert.Equal(t, 10, results[0].UnreadCount, "Group DM with 10 unread should be first")
	assert.Equal(t, 5, results[1].UnreadCount, "Channel with 5 unread should be second")
	assert.Equal(t, 3, results[2].UnreadCount, "DM with 3 unread should be third")

	mockClient.AssertExpectations(t)
}

func TestListUnreadChannelsOnly(t *testing.T) {
	mockClient := new(MockClient)

	channel1 := slack.Channel{}
	channel1.ID = "C1"
	channel1.Name = "general"
	channel1.IsChannel = true
	channel1.UnreadCount = 5
	channel1.UnreadCountDisplay = 5

	dm1 := slack.Channel{}
	dm1.ID = "D1"
	dm1.IsIM = true
	dm1.UnreadCount = 3
	dm1.UnreadCountDisplay = 3
	dm1.User = "U123"

	channels := []slack.Channel{channel1, dm1}

	mockClient.On("GetConversations", mock.Anything).Return(channels, "", nil)

	svc := NewUnreadService(mockClient)
	results, err := svc.ListUnread(UnreadOptions{ChannelsOnly: true})

	assert.NoError(t, err)
	assert.Len(t, results, 1, "Should return only channels")
	assert.Equal(t, "C1", results[0].ID)
	assert.True(t, results[0].IsChannel)

	mockClient.AssertExpectations(t)
}

func TestListUnreadDMsOnly(t *testing.T) {
	mockClient := new(MockClient)

	channel1 := slack.Channel{}
	channel1.ID = "C1"
	channel1.Name = "general"
	channel1.IsChannel = true
	channel1.UnreadCount = 5
	channel1.UnreadCountDisplay = 5

	dm1 := slack.Channel{}
	dm1.ID = "D1"
	dm1.IsIM = true
	dm1.UnreadCount = 3
	dm1.UnreadCountDisplay = 3
	dm1.User = "U123"

	groupDM := slack.Channel{}
	groupDM.ID = "G1"
	groupDM.Name = "mpdm-group"
	groupDM.IsMpIM = true
	groupDM.UnreadCount = 2
	groupDM.UnreadCountDisplay = 2

	channels := []slack.Channel{channel1, dm1, groupDM}

	mockClient.On("GetConversations", mock.Anything).Return(channels, "", nil)
	mockClient.On("GetUserInfo", "U123").Return(&slack.User{
		ID:   "U123",
		Name: "alice",
	}, nil)

	svc := NewUnreadService(mockClient)
	results, err := svc.ListUnread(UnreadOptions{DMsOnly: true})

	assert.NoError(t, err)
	assert.Len(t, results, 2, "Should return only DMs (1-on-1 and groups)")
	assert.True(t, results[0].IsIM || results[0].IsMpIM)
	assert.True(t, results[1].IsIM || results[1].IsMpIM)

	mockClient.AssertExpectations(t)
}

func TestListUnreadWithMinCount(t *testing.T) {
	mockClient := new(MockClient)

	channel1 := slack.Channel{}
	channel1.ID = "C1"
	channel1.Name = "important"
	channel1.IsChannel = true
	channel1.UnreadCount = 10
	channel1.UnreadCountDisplay = 10

	channel2 := slack.Channel{}
	channel2.ID = "C2"
	channel2.Name = "general"
	channel2.IsChannel = true
	channel2.UnreadCount = 2
	channel2.UnreadCountDisplay = 2

	channels := []slack.Channel{channel1, channel2}

	mockClient.On("GetConversations", mock.Anything).Return(channels, "", nil)

	svc := NewUnreadService(mockClient)
	results, err := svc.ListUnread(UnreadOptions{MinUnreadCount: 5})

	assert.NoError(t, err)
	assert.Len(t, results, 1, "Should return only conversations with >= 5 unread")
	assert.Equal(t, "C1", results[0].ID)
	assert.Equal(t, 10, results[0].UnreadCount)

	mockClient.AssertExpectations(t)
}

func TestListUnreadEmpty(t *testing.T) {
	mockClient := new(MockClient)

	// All channels have 0 unread
	channel1 := slack.Channel{}
	channel1.ID = "C1"
	channel1.Name = "general"
	channel1.IsChannel = true
	channel1.UnreadCount = 0

	channels := []slack.Channel{channel1}

	mockClient.On("GetConversations", mock.Anything).Return(channels, "", nil)

	svc := NewUnreadService(mockClient)
	results, err := svc.ListUnread(UnreadOptions{})

	assert.NoError(t, err)
	assert.Len(t, results, 0, "Should return empty list when nothing is unread")

	mockClient.AssertExpectations(t)
}

func TestListUnreadOrderByOldest(t *testing.T) {
	mockClient := new(MockClient)

	// Channels with different last_read timestamps
	channel1 := slack.Channel{}
	channel1.ID = "C1"
	channel1.Name = "newest"
	channel1.IsChannel = true
	channel1.UnreadCount = 5
	channel1.UnreadCountDisplay = 5
	channel1.LastRead = "1706123480.000000" // Most recent

	channel2 := slack.Channel{}
	channel2.ID = "C2"
	channel2.Name = "oldest"
	channel2.IsChannel = true
	channel2.UnreadCount = 3
	channel2.UnreadCountDisplay = 3
	channel2.LastRead = "1706123450.000000" // Oldest

	channel3 := slack.Channel{}
	channel3.ID = "C3"
	channel3.Name = "middle"
	channel3.IsChannel = true
	channel3.UnreadCount = 10
	channel3.UnreadCountDisplay = 10
	channel3.LastRead = "1706123465.000000" // Middle

	channels := []slack.Channel{channel1, channel2, channel3}

	mockClient.On("GetConversations", mock.Anything).Return(channels, "", nil)

	svc := NewUnreadService(mockClient)
	results, err := svc.ListUnread(UnreadOptions{OrderBy: "oldest"})

	assert.NoError(t, err)
	assert.Len(t, results, 3)

	// Should be ordered by LastRead ascending (oldest first)
	assert.Equal(t, "C2", results[0].ID, "Oldest unread should be first")
	assert.Equal(t, "1706123450.000000", results[0].LastRead)

	assert.Equal(t, "C3", results[1].ID, "Middle unread should be second")
	assert.Equal(t, "1706123465.000000", results[1].LastRead)

	assert.Equal(t, "C1", results[2].ID, "Newest unread should be last")
	assert.Equal(t, "1706123480.000000", results[2].LastRead)

	mockClient.AssertExpectations(t)
}

func TestListUnreadOrderByCount(t *testing.T) {
	mockClient := new(MockClient)

	channel1 := slack.Channel{}
	channel1.ID = "C1"
	channel1.Name = "low"
	channel1.IsChannel = true
	channel1.UnreadCount = 2
	channel1.UnreadCountDisplay = 2

	channel2 := slack.Channel{}
	channel2.ID = "C2"
	channel2.Name = "high"
	channel2.IsChannel = true
	channel2.UnreadCount = 15
	channel2.UnreadCountDisplay = 15

	channel3 := slack.Channel{}
	channel3.ID = "C3"
	channel3.Name = "medium"
	channel3.IsChannel = true
	channel3.UnreadCount = 7
	channel3.UnreadCountDisplay = 7

	channels := []slack.Channel{channel1, channel2, channel3}

	mockClient.On("GetConversations", mock.Anything).Return(channels, "", nil)

	svc := NewUnreadService(mockClient)
	// Default ordering should be by count
	results, err := svc.ListUnread(UnreadOptions{})

	assert.NoError(t, err)
	assert.Len(t, results, 3)

	// Should be ordered by UnreadCount descending (highest first)
	assert.Equal(t, "C2", results[0].ID, "Highest unread count should be first")
	assert.Equal(t, 15, results[0].UnreadCount)

	assert.Equal(t, "C3", results[1].ID, "Medium unread count should be second")
	assert.Equal(t, 7, results[1].UnreadCount)

	assert.Equal(t, "C1", results[2].ID, "Lowest unread count should be last")
	assert.Equal(t, 2, results[2].UnreadCount)

	mockClient.AssertExpectations(t)
}
