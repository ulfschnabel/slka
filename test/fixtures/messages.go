package fixtures

import "time"

// Message represents a Slack message for testing
type Message struct {
	Type      string
	User      string
	Text      string
	Timestamp string
	ThreadTS  string
}

// GetTestMessages returns test messages for a channel
func GetTestMessages(channelID string) []Message {
	now := time.Now().Unix()

	return []Message{
		{
			Type:      "message",
			User:      "U001",
			Text:      "Hello everyone!",
			Timestamp: formatMessageTimestamp(now - 3600),
		},
		{
			Type:      "message",
			User:      "U002",
			Text:      "Hi Alice!",
			Timestamp: formatMessageTimestamp(now - 3500),
		},
		{
			Type:      "message",
			User:      "U003",
			Text:      "Good morning team",
			Timestamp: formatMessageTimestamp(now - 3400),
		},
		{
			Type:      "message",
			User:      "U001",
			Text:      "Anyone up for lunch?",
			Timestamp: formatMessageTimestamp(now - 1800),
		},
		{
			Type:      "message",
			User:      "U002",
			Text:      "Sure! Where should we go?",
			Timestamp: formatMessageTimestamp(now - 1700),
		},
	}
}

func formatMessageTimestamp(unixTime int64) string {
	return time.Unix(unixTime, 0).Format("1234567890.000000")
}
