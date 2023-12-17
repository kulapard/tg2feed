// Package parser provides the parser for the telegram channel page.
package parser

import (
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"log"
	"net/http"
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
	return fmt.Sprintf("https://t.me/s/%s", chName)
}

// Parse returns the page object
func Parse(chName string) *Page {
	// Build web url
	channelURL := GetChannelWebURL(chName)

	// Request the HTML page.
	res, err := http.Get(channelURL) //nolint:gosec
	if err != nil {
		log.Fatal(err)
	}
	defer res.Body.Close()
	if res.StatusCode != 200 {
		log.Fatalf("status code error: %d %s", res.StatusCode, res.Status)
	}

	// Load the HTML document
	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		log.Fatal(err)
	}

	return GetPage(doc)
}
