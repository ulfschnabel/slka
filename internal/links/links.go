package links

import (
	"regexp"
)

// Link represents a parsed link
type Link struct {
	URL  string `json:"url"`
	Text string `json:"text,omitempty"`
}

var (
	// Slack link format: <url|text> or <url>
	slackLinkRegex = regexp.MustCompile(`<([^|>]+)(?:\|([^>]+))?>`)

	// Markdown link format: [text](url)
	markdownLinkRegex = regexp.MustCompile(`\[([^\]]+)\]\(([^\)]+)\)`)
)

// ParseSlackLinks extracts all links from Slack-formatted text
func ParseSlackLinks(text string) []Link {
	matches := slackLinkRegex.FindAllStringSubmatch(text, -1)
	links := make([]Link, 0, len(matches))

	for _, match := range matches {
		link := Link{
			URL: match[1],
		}
		if len(match) > 2 && match[2] != "" {
			link.Text = match[2]
		}
		links = append(links, link)
	}

	return links
}

// FormatLinksForSlack converts markdown links to Slack format
// Plain URLs and existing Slack links are passed through unchanged
func FormatLinksForSlack(text string) string {
	// Convert markdown [text](url) to Slack <url|text>
	result := markdownLinkRegex.ReplaceAllString(text, "<$2|$1>")
	return result
}
