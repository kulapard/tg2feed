package parser

import (
	"github.com/PuerkitoBio/goquery"
	"github.com/stretchr/testify/assert"
	"strings"
	"testing"
)

const testPageHTML = `
<html><body>
<i class="tgme_page_photo_image bgcolor0" data-content="нд">
	<img src="https://cdn4.cdn-telegram.org/file/img.jpg">
</i>
<div class="tgme_channel_info_header_title_wrap">
	<div class="tgme_channel_info_header_title"><span dir="auto">Some title</span></div>
	<div class="tgme_channel_info_header_labels"><i class="verified-icon"> ✔</i></div>
</div>
<div class="tgme_channel_info_header_username">
	<a href="https://t.me/telegram">@telegram</a>
</div>
<div class="tgme_channel_info_description">Some <b>page</b> description</div>
</body></html>
`

func TestGetPageTitle(t *testing.T) {
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(testPageHTML))
	assert.Nil(t, err)

	title := GetPageTitle(doc)
	assert.Equal(t, title, "Some title")
}

func TestGetPageLink(t *testing.T) {
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(testPageHTML))
	assert.Nil(t, err)

	link := GetPageLink(doc)
	assert.Equal(t, link, "https://t.me/telegram")

	// Empty link
	doc, err = goquery.NewDocumentFromReader(strings.NewReader(``))
	assert.Nil(t, err)

	link = GetPageLink(doc)
	assert.Equal(t, link, "")
}

func TestGetPageDescription(t *testing.T) {
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(testPageHTML))
	assert.Nil(t, err)

	description := GetPageDescriptionHTML(doc)
	assert.Equal(t, description, "Some <b>page</b> description")
}

func TestGetPageImageURL(t *testing.T) {
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(testPageHTML))
	assert.Nil(t, err)

	imgURL := GetPageImageURL(doc)
	assert.Equal(t, imgURL, "https://cdn4.cdn-telegram.org/file/img.jpg")

	// Empty image
	doc, err = goquery.NewDocumentFromReader(strings.NewReader(``))
	assert.Nil(t, err)

	imgURL = GetPageImageURL(doc)
	assert.Equal(t, imgURL, "")
}

func TestGetPage(t *testing.T) {
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(testPageHTML))
	assert.Nil(t, err)

	page := GetPage(doc)
	assert.Equal(t, page.Title, "Some title")
	assert.Equal(t, page.Link, "https://t.me/telegram")
	assert.Equal(t, page.Description, "Some <b>page</b> description")
	assert.Equal(t, page.ImageURL, "https://cdn4.cdn-telegram.org/file/img.jpg")
	assert.Nil(t, page.Posts)
}
