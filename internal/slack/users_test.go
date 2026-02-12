package slack

import (
	"testing"

	slackapi "github.com/slack-go/slack"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestListUsersSortedByUpdated(t *testing.T) {
	mockClient := new(MockClient)
	mockClient.On("GetUsers", mock.Anything).Return(
		[]slackapi.User{
			{ID: "U_OLD", Name: "old-user", RealName: "Old User", Updated: 1600000000},
			{ID: "U_NEW", Name: "new-user", RealName: "New User", Updated: 1706000000},
			{ID: "U_MID", Name: "mid-user", RealName: "Mid User", Updated: 1650000000},
		},
		nil,
	)

	svc := NewUserService(mockClient)
	result, err := svc.List(ListUsersOptions{})

	assert.NoError(t, err)
	assert.Len(t, result, 3)
	// Most recently updated first
	assert.Equal(t, "new-user", result[0].Name)
	assert.Equal(t, int64(1706000000), result[0].Updated)
	assert.Equal(t, "mid-user", result[1].Name)
	assert.Equal(t, int64(1650000000), result[1].Updated)
	assert.Equal(t, "old-user", result[2].Name)
	assert.Equal(t, int64(1600000000), result[2].Updated)
	mockClient.AssertExpectations(t)
}

func TestListUsersWithLimit(t *testing.T) {
	mockClient := new(MockClient)
	mockClient.On("GetUsers", mock.Anything).Return(
		[]slackapi.User{
			{ID: "U1", Name: "user1", Updated: 1706000000},
			{ID: "U2", Name: "user2", Updated: 1650000000},
			{ID: "U3", Name: "user3", Updated: 1600000000},
		},
		nil,
	)

	svc := NewUserService(mockClient)
	result, err := svc.List(ListUsersOptions{Limit: 2})

	assert.NoError(t, err)
	assert.Len(t, result, 2)
	// Sorted first, then limited — so we get the two most recently updated
	assert.Equal(t, "user1", result[0].Name)
	assert.Equal(t, "user2", result[1].Name)
	mockClient.AssertExpectations(t)
}
