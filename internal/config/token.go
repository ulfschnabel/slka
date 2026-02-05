package config

import "strings"

// TokenType represents the type of Slack token
type TokenType int

const (
	TokenTypeUnknown TokenType = iota
	TokenTypeBot                // xoxb-
	TokenTypeUser               // xoxp-
	TokenTypeApp                // xapp-
)

// DetectTokenType determines the type of a Slack token
func DetectTokenType(token string) TokenType {
	if token == "" {
		return TokenTypeUnknown
	}

	if strings.HasPrefix(token, "xoxb-") {
		return TokenTypeBot
	}
	if strings.HasPrefix(token, "xoxp-") {
		return TokenTypeUser
	}
	if strings.HasPrefix(token, "xapp-") {
		return TokenTypeApp
	}

	return TokenTypeUnknown
}

// IsUserToken checks if a token is a user token
func IsUserToken(token string) bool {
	return DetectTokenType(token) == TokenTypeUser
}

// IsBotToken checks if a token is a bot token
func IsBotToken(token string) bool {
	return DetectTokenType(token) == TokenTypeBot
}

// GetTokenTypeName returns a human-readable name for the token type
func GetTokenTypeName(tokenType TokenType) string {
	switch tokenType {
	case TokenTypeBot:
		return "Bot Token"
	case TokenTypeUser:
		return "User Token"
	case TokenTypeApp:
		return "App Token"
	default:
		return "Unknown Token"
	}
}

// GetTokenTypeDescription returns a description of what the token type is used for
func GetTokenTypeDescription(tokenType TokenType) string {
	switch tokenType {
	case TokenTypeBot:
		return "Bot tokens (xoxb-) allow your app to act as a bot user. Messages will appear from the bot."
	case TokenTypeUser:
		return "User tokens (xoxp-) allow your app to act as a specific user. Messages will appear from that user."
	case TokenTypeApp:
		return "App-level tokens (xapp-) are used for app-level features like Socket Mode."
	default:
		return "Unknown token type"
	}
}
