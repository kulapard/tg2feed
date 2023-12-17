// Package parser provides the parser for the telegram channel page.
package parser

import (
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

// getTitle returns the title of the selection.
func getTitle(s *goquery.Selection, _ ...string) (out interface{}, err error) {
	title, err := s.Html()
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
	title = ShortenText(title, 30)
	title = strings.TrimSpace(title)
	return title, nil
}

// getFormattedHTML returns the formatted html of the selection.
func getFormattedHTML(s *goquery.Selection, _ ...string) (out interface{}, err error) {
	formattedText := GetPostTextHTML(s)
	return formattedText, nil
}

// getImageURL returns the image url of the selection.
func getImageURL(s *goquery.Selection, _ ...string) (out interface{}, err error) {
	var imageURL string

	style, exists := s.Attr("style")
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

// getPostLink returns the post link.
func getPostLink(s *goquery.Selection, _ ...string) (out interface{}, err error) {
	baseURL, err := url.Parse("https://t.me/s/")
	if err != nil {
		return "", fmt.Errorf("invalid base url: %v error: %v", baseURL, err)
	}
	link, exists := s.Find(".tgme_widget_message").Attr("data-post")
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
func getGUID(s *goquery.Selection, _ ...string) (out interface{}, err error) {
	link, exists := s.Find(".tgme_widget_message").Attr("data-post")
	if !exists {
		return "", fmt.Errorf("data-post not found")
	}
	return GetGUID(link), nil
}

// getCreated returns the created date of the post.
func getCreated(s *goquery.Selection, _ ...string) (out interface{}, err error) {
	dt, exists := s.Attr("datetime")
	if !exists {
		return "", fmt.Errorf("datetime not found")
	}

	// Parse the datetime string
	return ParseDateTime(dt)
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
		Link    string    `pagser:"->getPostLink()"`
		ID      string    `pagser:"->getGUID()"`
		Created time.Time `pagser:".tgme_widget_message_date time->getCreated()"`
		Video   struct {
			URL string `pagser:"->attr(src)"`
		} `pagser:"video"`
		Images []struct {
			URL string `pagser:"->getImageURL()"`
		} `pagser:".tgme_widget_message_photo_wrap"`
		Previews []struct {
			Link        string `pagser:"->attr(href)"`
			ImageURL    string `pagser:".link_preview_image->getImageURL()"`
			VideoURL    string `pagser:"video->attr(src)"`
			SiteName    string `pagser:".link_preview_site_name->text()"`
			Title       string `pagser:".link_preview_title->text()"`
			Description string `pagser:".link_preview_description->getFormattedHTML()"`
		} `pagser:".tgme_widget_message_link_preview"`
	} `pagser:".tgme_widget_message_wrap"`
}

// Parse parses the page and returns the data model
func Parse(chName string) PageData {
	// Build web url
	channelURL := GetChannelWebURL(chName)

	// New default config
	p := pagser.New()

	// Register global functions
	p.RegisterFunc("getTitle", getTitle)
	p.RegisterFunc("getImageURL", getImageURL)
	p.RegisterFunc("getFormattedHTML", getFormattedHTML)
	p.RegisterFunc("getPostLink", getPostLink)
	p.RegisterFunc("getGUID", getGUID)
	p.RegisterFunc("getCreated", getCreated)

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
