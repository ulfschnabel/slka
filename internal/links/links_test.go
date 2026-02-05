package links

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParseSlackLinkWithText(t *testing.T) {
	// Parse Slack link format <url|text>
	text := "Check out <https://example.com|this link> for more"
	links := ParseSlackLinks(text)

	assert.Len(t, links, 1)
	assert.Equal(t, "https://example.com", links[0].URL)
	assert.Equal(t, "this link", links[0].Text)
}

func TestParseSlackLinkWithoutText(t *testing.T) {
	// Parse Slack link format <url> without display text
	text := "Visit <https://example.com>"
	links := ParseSlackLinks(text)

	assert.Len(t, links, 1)
	assert.Equal(t, "https://example.com", links[0].URL)
	assert.Empty(t, links[0].Text)
}

func TestParseMultipleLinks(t *testing.T) {
	text := "Visit <https://example.com|Example> and <https://test.com>"
	links := ParseSlackLinks(text)

	assert.Len(t, links, 2)
	assert.Equal(t, "https://example.com", links[0].URL)
	assert.Equal(t, "Example", links[0].Text)
	assert.Equal(t, "https://test.com", links[1].URL)
	assert.Empty(t, links[1].Text)
}

func TestParseNoLinks(t *testing.T) {
	text := "Just plain text"
	links := ParseSlackLinks(text)

	assert.Empty(t, links)
}

func TestFormatMarkdownLinkToSlack(t *testing.T) {
	// Convert markdown [text](url) to Slack <url|text>
	text := "Check out [this link](https://example.com) for more"
	result := FormatLinksForSlack(text)

	assert.Equal(t, "Check out <https://example.com|this link> for more", result)
}

func TestFormatMultipleMarkdownLinks(t *testing.T) {
	text := "Visit [Example](https://example.com) and [Test](https://test.com)"
	result := FormatLinksForSlack(text)

	assert.Equal(t, "Visit <https://example.com|Example> and <https://test.com|Test>", result)
}

func TestPassthroughSlackNativeLinks(t *testing.T) {
	// Don't double-encode already-formatted Slack links
	text := "Check out <https://example.com|this link> for more"
	result := FormatLinksForSlack(text)

	assert.Equal(t, text, result)
}

func TestPassthroughPlainURLs(t *testing.T) {
	// Plain URLs should pass through unchanged
	text := "Visit https://example.com for more"
	result := FormatLinksForSlack(text)

	assert.Equal(t, text, result)
}

func TestFormatMixedLinks(t *testing.T) {
	// Mix of markdown and plain URLs
	text := "Visit [Example](https://example.com) or https://test.com"
	result := FormatLinksForSlack(text)

	assert.Equal(t, "Visit <https://example.com|Example> or https://test.com", result)
}
