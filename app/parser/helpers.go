package parser

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"strings"
	"time"
	"unicode"
)

func GetSafeHTML(s *goquery.Selection) string {
	// Fix emoji
	s = FixEmoji(s)

	// Fix links
	s = FixLinks(s)

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

func ParseDateTime(dt string) (time.Time, error) {
	if dt == "" {
		return time.Now(), fmt.Errorf("can't parse empty datetime")
	}
	if ts, err := time.Parse(time.RFC3339, dt); err == nil {
		return ts, nil
	}
	return time.Now(), fmt.Errorf("can't parse datetime %s", dt)
}

func GetGUID(str string) string {
	hash := sha256.Sum256([]byte(str))
	hashStr := hex.EncodeToString(hash[:])
	return hashStr
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
