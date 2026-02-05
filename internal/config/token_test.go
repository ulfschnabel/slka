package config

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDetectTokenType(t *testing.T) {
	tests := []struct {
		name     string
		token    string
		expected TokenType
	}{
		{
			name:     "bot token",
			token:    "xoxb-1234567890-1234567890-abcdefghijklmnop",
			expected: TokenTypeBot,
		},
		{
			name:     "user token",
			token:    "xoxp-1234567890-1234567890-1234567890-abcdef",
			expected: TokenTypeUser,
		},
		{
			name:     "app token",
			token:    "xapp-1-A1234567890-1234567890-abcdef",
			expected: TokenTypeApp,
		},
		{
			name:     "empty token",
			token:    "",
			expected: TokenTypeUnknown,
		},
		{
			name:     "invalid token",
			token:    "invalid-token",
			expected: TokenTypeUnknown,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := DetectTokenType(tt.token)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestIsUserToken(t *testing.T) {
	assert.True(t, IsUserToken("xoxp-123-456-789-abc"))
	assert.False(t, IsUserToken("xoxb-123-456-abc"))
	assert.False(t, IsUserToken(""))
	assert.False(t, IsUserToken("invalid"))
}

func TestIsBotToken(t *testing.T) {
	assert.True(t, IsBotToken("xoxb-123-456-abc"))
	assert.False(t, IsBotToken("xoxp-123-456-789-abc"))
	assert.False(t, IsBotToken(""))
	assert.False(t, IsBotToken("invalid"))
}

func TestGetTokenTypeName(t *testing.T) {
	tests := []struct {
		tokenType TokenType
		expected  string
	}{
		{TokenTypeBot, "Bot Token"},
		{TokenTypeUser, "User Token"},
		{TokenTypeApp, "App Token"},
		{TokenTypeUnknown, "Unknown Token"},
	}

	for _, tt := range tests {
		t.Run(tt.expected, func(t *testing.T) {
			result := GetTokenTypeName(tt.tokenType)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestGetTokenTypeDescription(t *testing.T) {
	// Just verify they all return non-empty strings
	assert.NotEmpty(t, GetTokenTypeDescription(TokenTypeBot))
	assert.NotEmpty(t, GetTokenTypeDescription(TokenTypeUser))
	assert.NotEmpty(t, GetTokenTypeDescription(TokenTypeApp))
	assert.NotEmpty(t, GetTokenTypeDescription(TokenTypeUnknown))
}
