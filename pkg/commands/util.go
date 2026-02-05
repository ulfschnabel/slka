package commands

import (
	"fmt"
	"strconv"
	"time"
)

// parseTimestamp parses a timestamp string (Unix timestamp or ISO8601)
func parseTimestamp(s string) (int64, error) {
	// Try parsing as Unix timestamp
	if ts, err := strconv.ParseInt(s, 10, 64); err == nil {
		return ts, nil
	}

	// Try parsing as ISO8601
	layouts := []string{
		time.RFC3339,
		"2006-01-02T15:04:05",
		"2006-01-02",
	}

	for _, layout := range layouts {
		if t, err := time.Parse(layout, s); err == nil {
			return t.Unix(), nil
		}
	}

	return 0, fmt.Errorf("invalid timestamp format: %s", s)
}
