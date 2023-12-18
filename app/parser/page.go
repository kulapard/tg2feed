package parser

import "github.com/PuerkitoBio/goquery"

// Page represents a page from the telegram channel
type Page struct {
	Title       string
	Link        string
	Description string
	ImageURL    string
	Posts       []*Post
}

// GetPageTitle returns the page title
func GetPageTitle(doc *goquery.Document) string {
	return doc.Find(".tgme_channel_info_header_title").Text()
}

// GetPageLink returns the page link
func GetPageLink(doc *goquery.Document) string {
	if link, exists := doc.Find(".tgme_channel_info_header_username a").Attr("href"); exists {
		return link
	}
	return ""
}

// GetPageDescriptionHTML returns the page description html
func GetPageDescriptionHTML(doc *goquery.Document) string {
	return GetSafeHTML(doc.Find(".tgme_channel_info_description"))
}

// GetPageImageURL returns the page image url
func GetPageImageURL(doc *goquery.Document) string {
	if imgURL, exists := doc.Find(".tgme_page_photo_image img").Attr("src"); exists {
		return imgURL
	}
	return ""
}

// GetPage returns the page object
func GetPage(doc *goquery.Document) *Page {
	return &Page{
		Title:       GetPageTitle(doc),
		Link:        GetPageLink(doc),
		Description: GetPageDescriptionHTML(doc),
		ImageURL:    GetPageImageURL(doc),
		Posts:       GetPosts(doc),
	}
}
