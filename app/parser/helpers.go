package parser

import (
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"strings"
	"time"
	"unicode"
)

// GetSafeHTML returns the HTML string without unsafe tags
func GetSafeHTML(s *goquery.Selection) string {
	// Fix emoji
	s = FixEmoji(s)

	// Fix links
	s = FixLinks(s)

	// Remove unsafe tags
	s = RemoveUnsafeTags(s)

	html, err := s.Html()
	if err != nil {
		return ""
	}

	// Remove new lines
	html = strings.ReplaceAll(html, "\n", "")

	// Replace \t with space
	html = strings.ReplaceAll(html, "\t", " ")

	return strings.TrimSpace(html)
}

// ParseDateTime parses the datetime string and returns the time.Time object
func ParseDateTime(dt string) (time.Time, error) {
	if dt == "" {
		return time.Now(), fmt.Errorf("can't parse empty datetime")
	}
	if ts, err := time.Parse(time.RFC3339, dt); err == nil {
		return ts, nil
	}
	return time.Now(), fmt.Errorf("can't parse datetime %s", dt)
}

// ShortenText shortens the text to the specified length.
func ShortenText(text string, maxLength int) string {
	lastSpaceIx := -1
	length := 0

	// Iterate over runes to get correct index
	for i, r := range text {
		if unicode.IsSpace(r) {
			lastSpaceIx = i
		}
		length++
		if length > maxLength {
			if lastSpaceIx != -1 {
				// String is longer than max and has a space.
				// Let's cut it off by space to avoid cutting words in half
				return text[:lastSpaceIx] + "..."
			}
			// String is longer than max and has no space. Let's cut it off by force
			return text[:i] + "..."
		}
	}
	// String is already shorter than max
	return text
}
