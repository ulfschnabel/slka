package slack

import (
	"errors"
	"testing"

	"github.com/slack-go/slack"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestAddReactionSuccess(t *testing.T) {
	mockClient := new(MockClient)
	mockClient.On("AddReaction", "thumbsup", mock.MatchedBy(func(item slack.ItemRef) bool {
		return item.Channel == "C123" && item.Timestamp == "1706123456.789000"
	})).Return(nil)

	svc := NewReactionService(mockClient)
	err := svc.AddReaction("C123", "1706123456.789000", "thumbsup")

	assert.NoError(t, err)
	mockClient.AssertExpectations(t)
}

func TestAddReactionStripsColons(t *testing.T) {
	mockClient := new(MockClient)
	mockClient.On("AddReaction", "thumbsup", mock.Anything).Return(nil)

	svc := NewReactionService(mockClient)
	err := svc.AddReaction("C123", "1706123456.789000", ":thumbsup:")

	assert.NoError(t, err)
	mockClient.AssertExpectations(t)
}

func TestAddReactionAPIError(t *testing.T) {
	mockClient := new(MockClient)
	mockClient.On("AddReaction", "thumbsup", mock.Anything).Return(errors.New("already_reacted"))

	svc := NewReactionService(mockClient)
	err := svc.AddReaction("C123", "1706123456.789000", "thumbsup")

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "failed to add reaction")
	mockClient.AssertExpectations(t)
}

func TestRemoveReactionSuccess(t *testing.T) {
	mockClient := new(MockClient)
	mockClient.On("RemoveReaction", "eyes", mock.MatchedBy(func(item slack.ItemRef) bool {
		return item.Channel == "C123" && item.Timestamp == "1706123456.789000"
	})).Return(nil)

	svc := NewReactionService(mockClient)
	err := svc.RemoveReaction("C123", "1706123456.789000", "eyes")

	assert.NoError(t, err)
	mockClient.AssertExpectations(t)
}

func TestRemoveReactionStripsColons(t *testing.T) {
	mockClient := new(MockClient)
	mockClient.On("RemoveReaction", "eyes", mock.Anything).Return(nil)

	svc := NewReactionService(mockClient)
	err := svc.RemoveReaction("C123", "1706123456.789000", ":eyes:")

	assert.NoError(t, err)
	mockClient.AssertExpectations(t)
}

func TestListReactionsSuccess(t *testing.T) {
	mockClient := new(MockClient)

	// Mock GetReactions
	mockClient.On("GetReactions", mock.Anything, mock.Anything).Return(
		[]slack.ItemReaction{
			{
				Name:  "thumbsup",
				Count: 2,
				Users: []string{"U123", "U456"},
			},
			{
				Name:  "eyes",
				Count: 1,
				Users: []string{"U789"},
			},
		},
		nil,
	)

	// Mock GetConversationHistory to get message text
	mockClient.On("GetConversationHistory", mock.MatchedBy(func(params *slack.GetConversationHistoryParameters) bool {
		return params.ChannelID == "C123" && params.Latest == "1706123456.789000"
	})).Return(
		&slack.GetConversationHistoryResponse{
			Messages: []slack.Message{
				{
					Msg: slack.Msg{
						Text:      "Hello team!",
						Timestamp: "1706123456.789000",
					},
				},
			},
		},
		nil,
	)

	svc := NewReactionService(mockClient)
	result, err := svc.ListReactions("C123", "1706123456.789000")

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, "1706123456.789000", result.Timestamp)
	assert.Equal(t, "C123", result.Channel)
	assert.Equal(t, "Hello team!", result.MessageText)
	assert.Len(t, result.Reactions, 2)
	assert.Equal(t, 3, result.TotalReactionCount)
	mockClient.AssertExpectations(t)
}

func TestListReactionsNoReactions(t *testing.T) {
	mockClient := new(MockClient)

	mockClient.On("GetReactions", mock.Anything, mock.Anything).Return(
		[]slack.ItemReaction{},
		nil,
	)

	mockClient.On("GetConversationHistory", mock.Anything).Return(
		&slack.GetConversationHistoryResponse{
			Messages: []slack.Message{
				{
					Msg: slack.Msg{
						Text:      "No reactions here",
						Timestamp: "1706123456.789000",
					},
				},
			},
		},
		nil,
	)

	svc := NewReactionService(mockClient)
	result, err := svc.ListReactions("C123", "1706123456.789000")

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Len(t, result.Reactions, 0)
	assert.Equal(t, 0, result.TotalReactionCount)
	mockClient.AssertExpectations(t)
}

func TestCheckAcknowledgmentWithReactions(t *testing.T) {
	mockClient := new(MockClient)

	// Mock GetConversationHistory to get message
	mockClient.On("GetConversationHistory", mock.Anything).Return(
		&slack.GetConversationHistoryResponse{
			Messages: []slack.Message{
				{
					Msg: slack.Msg{
						Timestamp:  "1706123456.789000",
						User:       "U123", // Message author
						ReplyCount: 0,
						Reactions: []slack.ItemReaction{
							{
								Name:  "thumbsup",
								Count: 2,
								Users: []string{"U456", "U789"}, // Others reacted
							},
						},
					},
				},
			},
		},
		nil,
	)

	svc := NewReactionService(mockClient)
	result, err := svc.CheckAcknowledgment("C123", "1706123456.789000")

	assert.NoError(t, err)
	assert.True(t, result.IsAcknowledged)
	assert.True(t, result.HasReactions)
	assert.False(t, result.HasReplies)
	assert.Equal(t, 2, result.ReactionCount)
	assert.Equal(t, 0, result.ReplyCount)
	assert.Equal(t, "U123", result.MessageAuthor)
	assert.ElementsMatch(t, []string{"U456", "U789"}, result.ReactedUsers)
	mockClient.AssertExpectations(t)
}

func TestCheckAcknowledgmentWithReplies(t *testing.T) {
	mockClient := new(MockClient)

	// Mock GetConversationHistory
	mockClient.On("GetConversationHistory", mock.Anything).Return(
		&slack.GetConversationHistoryResponse{
			Messages: []slack.Message{
				{
					Msg: slack.Msg{
						Timestamp:  "1706123456.789000",
						User:       "U123",
						ReplyCount: 2,
					},
				},
			},
		},
		nil,
	)

	// Mock GetConversationReplies
	mockClient.On("GetConversationReplies", mock.MatchedBy(func(params *slack.GetConversationRepliesParameters) bool {
		return params.ChannelID == "C123" && params.Timestamp == "1706123456.789000"
	})).Return(
		[]slack.Message{
			{
				Msg: slack.Msg{
					Timestamp: "1706123456.789000",
					User:      "U123", // Parent message
				},
			},
			{
				Msg: slack.Msg{
					Timestamp: "1706123457.000000",
					User:      "U456", // Reply from someone else
				},
			},
			{
				Msg: slack.Msg{
					Timestamp: "1706123458.000000",
					User:      "U789", // Another reply
				},
			},
		},
		false, // hasMore
		"",    // nextCursor
		nil,   // error
	)

	svc := NewReactionService(mockClient)
	result, err := svc.CheckAcknowledgment("C123", "1706123456.789000")

	assert.NoError(t, err)
	assert.True(t, result.IsAcknowledged)
	assert.False(t, result.HasReactions)
	assert.True(t, result.HasReplies)
	assert.Equal(t, 0, result.ReactionCount)
	assert.Equal(t, 2, result.ReplyCount) // Excludes parent message
	mockClient.AssertExpectations(t)
}

func TestCheckAcknowledgmentWithBoth(t *testing.T) {
	mockClient := new(MockClient)

	mockClient.On("GetConversationHistory", mock.Anything).Return(
		&slack.GetConversationHistoryResponse{
			Messages: []slack.Message{
				{
					Msg: slack.Msg{
						Timestamp:  "1706123456.789000",
						User:       "U123",
						ReplyCount: 1,
						Reactions: []slack.ItemReaction{
							{
								Name:  "eyes",
								Count: 1,
								Users: []string{"U456"},
							},
						},
					},
				},
			},
		},
		nil,
	)

	mockClient.On("GetConversationReplies", mock.Anything).Return(
		[]slack.Message{
			{
				Msg: slack.Msg{
					Timestamp: "1706123456.789000",
					User:      "U123",
				},
			},
			{
				Msg: slack.Msg{
					Timestamp: "1706123457.000000",
					User:      "U789",
				},
			},
		},
		false, // hasMore
		"",    // nextCursor
		nil,   // error
	)

	svc := NewReactionService(mockClient)
	result, err := svc.CheckAcknowledgment("C123", "1706123456.789000")

	assert.NoError(t, err)
	assert.True(t, result.IsAcknowledged)
	assert.True(t, result.HasReactions)
	assert.True(t, result.HasReplies)
	assert.Equal(t, 1, result.ReactionCount)
	assert.Equal(t, 1, result.ReplyCount)
	mockClient.AssertExpectations(t)
}

func TestCheckAcknowledgmentSelfOnly(t *testing.T) {
	mockClient := new(MockClient)

	mockClient.On("GetConversationHistory", mock.Anything).Return(
		&slack.GetConversationHistoryResponse{
			Messages: []slack.Message{
				{
					Msg: slack.Msg{
						Timestamp:  "1706123456.789000",
						User:       "U123",
						ReplyCount: 0,
						Reactions: []slack.ItemReaction{
							{
								Name:  "thumbsup",
								Count: 1,
								Users: []string{"U123"}, // Only author reacted
							},
						},
					},
				},
			},
		},
		nil,
	)

	svc := NewReactionService(mockClient)
	result, err := svc.CheckAcknowledgment("C123", "1706123456.789000")

	assert.NoError(t, err)
	assert.False(t, result.IsAcknowledged) // Self-reaction doesn't count
	assert.True(t, result.HasReactions)
	assert.False(t, result.HasReplies)
	assert.Equal(t, 0, result.ReactionCount) // Count excludes self
	assert.Len(t, result.ReactedUsers, 0)
	mockClient.AssertExpectations(t)
}

func TestCheckAcknowledgmentNone(t *testing.T) {
	mockClient := new(MockClient)

	mockClient.On("GetConversationHistory", mock.Anything).Return(
		&slack.GetConversationHistoryResponse{
			Messages: []slack.Message{
				{
					Msg: slack.Msg{
						Timestamp:  "1706123456.789000",
						User:       "U123",
						ReplyCount: 0,
					},
				},
			},
		},
		nil,
	)

	svc := NewReactionService(mockClient)
	result, err := svc.CheckAcknowledgment("C123", "1706123456.789000")

	assert.NoError(t, err)
	assert.False(t, result.IsAcknowledged)
	assert.False(t, result.HasReactions)
	assert.False(t, result.HasReplies)
	assert.Equal(t, 0, result.ReactionCount)
	assert.Equal(t, 0, result.ReplyCount)
	mockClient.AssertExpectations(t)
}
