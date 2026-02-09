package fixtures

import "time"

// Channel represents a Slack channel for testing
type Channel struct {
	ID                 string
	Name               string
	IsChannel          bool
	IsPrivate          bool
	IsIM               bool
	IsMpIM             bool
	UnreadCount        int
	UnreadCountDisplay int
	LastRead           string
	User               string // For IMs
	NumMembers         int
}

// GetTestChannels returns a set of test channels with various states
func GetTestChannels() []Channel {
	now := time.Now().Unix()
	oldTimestamp := now - 86400 // 1 day ago

	return []Channel{
		{
			ID:                 "C001",
			Name:               "general",
			IsChannel:          true,
			IsPrivate:          false,
			UnreadCount:        5,
			UnreadCountDisplay: 5,
			LastRead:           formatTimestamp(oldTimestamp),
			NumMembers:         50,
		},
		{
			ID:                 "C002",
			Name:               "engineering",
			IsChannel:          true,
			IsPrivate:          false,
			UnreadCount:        12,
			UnreadCountDisplay: 12,
			LastRead:           formatTimestamp(oldTimestamp - 3600),
			NumMembers:         25,
		},
		{
			ID:                 "C003",
			Name:               "random",
			IsChannel:          true,
			IsPrivate:          false,
			UnreadCount:        0,
			UnreadCountDisplay: 0,
			LastRead:           formatTimestamp(now),
			NumMembers:         100,
		},
		{
			ID:                 "C004",
			Name:               "secret-project",
			IsChannel:          true,
			IsPrivate:          true,
			UnreadCount:        3,
			UnreadCountDisplay: 3,
			LastRead:           formatTimestamp(oldTimestamp - 7200),
			NumMembers:         5,
		},
		{
			ID:        "D001",
			IsIM:      true,
			User:      "U001",
			UnreadCount: 2,
			UnreadCountDisplay: 2,
			LastRead:  formatTimestamp(oldTimestamp - 1800),
		},
		{
			ID:        "D002",
			IsIM:      true,
			User:      "U002",
			UnreadCount: 0,
			UnreadCountDisplay: 0,
			LastRead:  formatTimestamp(now),
		},
		{
			ID:                 "G001",
			Name:               "mpdm-alice--bob--charlie-1",
			IsMpIM:             true,
			UnreadCount:        8,
			UnreadCountDisplay: 8,
			LastRead:           formatTimestamp(oldTimestamp - 5400),
		},
	}
}

func formatTimestamp(unixTime int64) string {
	return time.Unix(unixTime, 0).Format("1234567890.000000")
}
