// Package parser provides the parser for the telegram channel page.
package parser

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"github.com/foolin/pagser"
	"io"
	"log"
	"net/http"
	"net/url"
	"regexp"
	"strings"
	"time"
	"unicode"
)

func getRawHTML(pageURL string) (string, error) {
	// Request the HTML page.
	res, err := http.Get(pageURL) //nolint:gosec // url id build from trusted source
	if err != nil {
		log.Fatal(err)
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return "", fmt.Errorf("status code error: %d %s", res.StatusCode, res.Status)
	}

	// Read the response body
	html, err := io.ReadAll(res.Body)
	if err != nil {
		return "", err
	}
	return string(html), nil
}

// shortenText shortens the text to the specified length.
func shortenText(text string, maxLength int) string {
	lastSpaceIx := -1
	length := 0

	// Iterate over runes to get correct index
	for i, r := range text {
		if unicode.IsSpace(r) {
			lastSpaceIx = i
		}
		length++
		if length >= maxLength {
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

// getTitle returns the title of the selection.
func getTitle(node *goquery.Selection, _ ...string) (out interface{}, err error) {
	title, err := node.Html()
	if err != nil {
		return "", fmt.Errorf("title not found")
	}
	// leave text until the first <br>
	brIndex := strings.Index(title, "<br")
	if brIndex != -1 {
		title = title[:brIndex]
	}

	// remove all tags using regexp
	title = regexp.MustCompile("<[^>]*>").ReplaceAllString(title, "")
	title = shortenText(title, 30)
	title = strings.TrimSpace(title)
	return title, nil
}

// getFormattedHTML returns the formatted html of the selection.
func getFormattedHTML(node *goquery.Selection, _ ...string) (out interface{}, err error) {
	var paragraphs []string
	html, err := node.Html()
	if err != nil {
		return "", err
	}
	node.SetHtml("")

	// split by <br>
	parts := strings.Split(html, "<br/>")
	if len(parts) <= 1 {
		parts = strings.Split(html, "<br>")
	}
	for _, part := range parts {
		// remove zero-width space
		part = strings.ReplaceAll(part, "\u200b", "")

		// ignore empty paragraphs
		if len(part) > 0 {
			paragraphs = append(paragraphs, fmt.Sprintf("<p>%s</p>", strings.TrimSpace(part)))
		}
	}
	formattedText := strings.Join(paragraphs, "\n")

	// replace <tg-spoiler>text</tg-spoiler> with <spoiler>text</spoiler>
	formattedText = strings.ReplaceAll(formattedText, "<tg-spoiler>", "<spoiler>")
	formattedText = strings.ReplaceAll(formattedText, "</tg-spoiler>", "</spoiler>")

	// replace https://t.me/... with https://t.me/s/...
	formattedText = strings.ReplaceAll(formattedText, "https://t.me/", "https://t.me/s/")

	return formattedText, nil
}

// getImageURL returns the image url of the selection.
func getImageURL(node *goquery.Selection, _ ...string) (out interface{}, err error) {
	var imageURL string

	style, exists := node.Attr("style")
	if exists {
		// Extract URL from the style attribute
		const stylePrefix = "background-image:url('"
		const styleSuffix = "')"
		start := strings.Index(style, stylePrefix)
		if start != -1 {
			start += len(stylePrefix)
			end := strings.Index(style[start:], styleSuffix)
			if end != -1 {
				imageURL = style[start : start+end]
			}
		}
	}
	return imageURL, nil
}

// getVideoURL returns the video url of the selection.
func getVideoURL(node *goquery.Selection, _ ...string) (out interface{}, err error) {
	videoURL, _ := node.Find("video").Attr("src")
	return videoURL, nil
}

// getPostLink returns the post link.
func getPostLink(node *goquery.Selection, args ...string) (out interface{}, err error) {
	if len(args) < 1 {
		return "", fmt.Errorf("args must has baseURL")
	}
	baseURL, err := url.Parse(args[0])
	if err != nil {
		return "", fmt.Errorf("invalid base url: %v error: %v", baseURL, err)
	}
	link, exists := node.Find(".tgme_widget_message").Attr("data-post")
	if !exists {
		return "", fmt.Errorf("data-post not found")
	}
	hrefURL, err := url.Parse(link)
	if err != nil {
		return "", err
	}
	return baseURL.ResolveReference(hrefURL), nil
}

// getGUID returns the guid of the selection based on the post link. It uses the sha256 hash of the link.
func getGUID(node *goquery.Selection, _ ...string) (out interface{}, err error) {
	link, exists := node.Find(".tgme_widget_message").Attr("data-post")
	if !exists {
		return "", fmt.Errorf("data-post not found")
	}
	hash := sha256.Sum256([]byte(link))
	hashStr := hex.EncodeToString(hash[:])
	return hashStr, nil
}

// getCreated returns the created date of the post.
func getCreated(node *goquery.Selection, _ ...string) (out interface{}, err error) {
	datetimeStr, exists := node.Attr("datetime")
	if !exists {
		return "", fmt.Errorf("datetime not found")
	}

	// Parse the datetime string
	const layout = "2006-01-02T15:04:05Z07:00"
	parsedTime, err := time.Parse(layout, datetimeStr)
	if err != nil {
		log.Fatalf("Error parsing time: %v", err)
	}
	return parsedTime, nil
}

// PageData is the data model for the page
type PageData struct {
	Title       string `pagser:".tgme_channel_info_header_title"`
	Link        string `pagser:".tgme_channel_info_header_username a->attr(href)"`
	Description string `pagser:".tgme_channel_info_description->getFormattedHTML()"`
	ImageURL    string `pagser:".tgme_page_photo_image img->attr(src)"`
	Posts       []*struct {
		Title   string    `pagser:".tgme_widget_message_text->getTitle()"`
		Text    string    `pagser:".tgme_widget_message_text->getFormattedHTML()"`
		Link    string    `pagser:"->getPostLink('https://t.me/s/')"`
		ID      string    `pagser:"->getGUID()"`
		Created time.Time `pagser:".tgme_widget_message_date time->getCreated()"`
		Video   struct {
			URL string `pagser:"->getVideoURL()"`
		} `pagser:".tgme_widget_message_video_player"`
		Images []struct {
			URL string `pagser:"->getImageURL()"`
		} `pagser:".tgme_widget_message_photo_wrap"`
		Previews []struct {
			Link        string `pagser:"->attr(href)"`
			ImageURL    string `pagser:".link_preview_image->getImageURL()"`
			VideoURL    string `pagser:".link_preview_video_wrap->getVideoURL()"`
			SiteName    string `pagser:".link_preview_site_name->text()"`
			Title       string `pagser:".link_preview_title->text()"`
			Description string `pagser:".link_preview_description->getFormattedHTML()"`
		} `pagser:".tgme_widget_message_link_preview"`
	} `pagser:".tgme_widget_message_wrap"`
}

// Parse parses the page and returns the data model
func Parse(channel string) PageData {
	// Remove @ from channel name
	channel = strings.ReplaceAll(channel, "@", "")

	// Remove spaces from channel name
	channel = strings.ReplaceAll(channel, " ", "")

	// Build channel url
	channelURL := fmt.Sprintf("https://t.me/s/%s", channel)

	// New default config
	p := pagser.New()

	// Register global functions
	p.RegisterFunc("getTitle", getTitle)
	p.RegisterFunc("getImageURL", getImageURL)
	p.RegisterFunc("getFormattedHTML", getFormattedHTML)
	p.RegisterFunc("getPostLink", getPostLink)
	p.RegisterFunc("getGUID", getGUID)
	p.RegisterFunc("getCreated", getCreated)
	p.RegisterFunc("getVideoURL", getVideoURL)

	// data parser model
	var data PageData

	// load html
	rawPageHTML, err := getRawHTML(channelURL)
	if err != nil {
		log.Fatal(err)
	}

	// parse data
	err = p.Parse(&data, rawPageHTML)
	if err != nil {
		log.Fatal(err)
	}

	return data
}
