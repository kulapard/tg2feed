// Package parser provides the parser for the telegram channel page.
package parser

import (
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"log"
	"net/http"
	"net/url"
	"strings"
)

// GetChannelWebURL returns the channel web url based on the channel name
func GetChannelWebURL(chName string) string {
	// Remove @ from chName name
	chName = strings.ReplaceAll(chName, "@", "")

	// Remove prefix https://t.me/s/ from channel name
	chName = strings.ReplaceAll(chName, "https://t.me/s/", "")

	// Remove prefix https://t.me/ from channel name
	chName = strings.ReplaceAll(chName, "https://t.me/", "")

	// Remove spaces from chName name
	chName = strings.TrimSpace(chName)

	// Build web url
	baseURL, _ := url.Parse("https://t.me/s/")
	if chPath, err := url.Parse(chName); err == nil {
		return baseURL.JoinPath(chPath.Path).String()
	}
	return ""
}

// Parse returns the page object
func Parse(chName string) (*Page, error) {
	// Build web url
	channelURL := GetChannelWebURL(chName)

	// Request the HTML page.
	res, err := http.Get(channelURL) //nolint:gosec // tolerate security risk
	if err != nil {
		log.Fatal(err)
	}
	defer res.Body.Close()
	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("status code error: %s", res.Status)
	}

	// Load the HTML document
	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		return nil, fmt.Errorf("can't parse HTML: %w", err)
	}
	return GetPage(doc), nil
}
