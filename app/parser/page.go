package parser

import "github.com/PuerkitoBio/goquery"

type Page struct {
	Title       string
	Link        string
	Description string
	ImageURL    string
	Posts       []*Post
}

func GetPageTitle(doc *goquery.Document) string {
	return doc.Find(".tgme_channel_info_header_title").Text()
}

func GetPageLink(doc *goquery.Document) string {
	link, exists := doc.Find(".tgme_channel_info_header_username a").Attr("href")
	if !exists {
		link = ""
	}
	return link
}

func GetPageDescriptionHtml(doc *goquery.Document) string {
	return GetSafeHTML(doc.Find(".tgme_channel_info_description"))
}
func GetPageImageURL(doc *goquery.Document) string {
	imgURL, exists := doc.Find(".tgme_page_photo_image img").Attr("src")
	if !exists {
		imgURL = ""
	}
	return imgURL
}

func GetPage(doc *goquery.Document) *Page {
	return &Page{
		Title:       GetPageTitle(doc),
		Link:        GetPageLink(doc),
		Description: GetPageDescriptionHtml(doc),
		ImageURL:    GetPageImageURL(doc),
		Posts:       GetPosts(doc),
	}
}
