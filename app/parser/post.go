package parser

import (
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"log"
	"net/url"
	"regexp"
	"strings"
	"time"
)

// Post represents a post from the telegram channel
type Post struct {
	Title   string
	Text    string
	Link    string
	ID      string
	Created time.Time
	Videos  []string
	Images  []string
}

// GetPosts returns all posts from the page
func GetPosts(doc *goquery.Document) []*Post {
	var posts []*Post

	doc.Find(".tgme_widget_message_wrap").Each(func(i int, s *goquery.Selection) {
		postLink := GetPostLink(s)
		text := GetPostTextHTML(s)
		posts = append(posts, &Post{
			Title:   GetPostTitle(text),
			Text:    text,
			Link:    postLink,
			Created: GetPostCreated(s),
			Videos:  GetVideos(s),
			Images:  GetImages(s),
		})
	})
	return posts
}

// GetPostTitle returns the post title
func GetPostTitle(text string) string {
	// Get only string before the first <br>
	title := strings.Split(text, "<br")[0]

	// Remove all tags using regexp
	title = regexp.MustCompile("<[^>]*>").ReplaceAllString(title, "")

	// Remove all new lines
	title = strings.TrimSpace(title)

	// Shorten the title to 30 characters
	return ShortenText(title, 30)
}

// GetPostTextHTML returns the post text as HTML
func GetPostTextHTML(s *goquery.Selection) string {
	s = s.Find(".tgme_widget_message_text")

	html := GetSafeHTML(s)

	var paragraphs []string

	// split by <br>
	parts := strings.Split(html, "<br/>")
	if len(parts) <= 1 {
		parts = strings.Split(html, "<br>")
	}
	for _, part := range parts {
		// ignore empty paragraphs
		if len(part) > 0 {
			paragraphs = append(paragraphs, fmt.Sprintf("<p>%s</p>", strings.TrimSpace(part)))
		}
	}
	formattedText := strings.Join(paragraphs, "\n")

	// replace https://t.me/... with https://t.me/s/...
	formattedText = strings.ReplaceAll(formattedText, "https://t.me/", "https://t.me/s/")

	return formattedText

}

// GetPostLink returns the post link
func GetPostLink(s *goquery.Selection) string {
	baseURL, err := url.Parse("https://t.me/s/")
	if err != nil {
		return ""
	}
	link, exists := s.Find(".tgme_widget_message").Attr("data-post")
	if !exists {
		return ""
	}
	hrefURL, err := url.Parse(link)
	if err != nil {
		return ""
	}
	return baseURL.ResolveReference(hrefURL).String()
}

// GetPostCreated returns the post created datetime
func GetPostCreated(s *goquery.Selection) time.Time {
	created, exists := s.Find(".tgme_widget_message_date time").Attr("datetime")
	if !exists {
		log.Print("[ERROR] post created datetime not found, fallback with time.Now()")
		return time.Now()
	}
	ts, err := ParseDateTime(created)
	if err != nil {
		log.Printf("[ERROR] failed to parse created: %v, fallback with time.Now()", err)
		return time.Now()
	}
	return ts
}

// GetVideos returns all videos from the post
func GetVideos(s *goquery.Selection) []string {
	var videos []string
	s.Find("video").Each(func(i int, s *goquery.Selection) {
		videoURL, exists := s.Attr("src")
		if exists {
			videos = append(videos, videoURL)
		}
	})
	return videos
}

// extractImageURLFromStyle extracts image URL from the style attribute
func extractImageURLFromStyle(s *goquery.Selection) string {
	style, exists := s.Attr("style")
	if !exists {
		return ""
	}
	// Extract URL from the style attribute
	const stylePrefix = "background-image:url('"
	const styleSuffix = "')"
	start := strings.Index(style, stylePrefix)
	if start != -1 {
		start += len(stylePrefix)
		end := strings.Index(style[start:], styleSuffix)
		if end != -1 {
			imageURL := style[start : start+end]
			return imageURL
		}
	}
	return ""
}

// GetImages returns all images from the post
func GetImages(s *goquery.Selection) []string {
	var images []string
	s.Find(".tgme_widget_message_photo_wrap").Each(func(i int, s *goquery.Selection) {
		if imageURL := extractImageURLFromStyle(s); imageURL != "" {
			images = append(images, imageURL)
		}
	})
	s.Find("img").Each(func(i int, s *goquery.Selection) {
		if imageURL, exists := s.Attr("src"); exists {
			images = append(images, imageURL)
		}
	})
	return images
}
