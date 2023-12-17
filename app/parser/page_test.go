package parser

import (
	"github.com/PuerkitoBio/goquery"
	"github.com/stretchr/testify/assert"
	"strings"
	"testing"
)

const testPageHTML = `
<html><body>
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
}

func TestGetPageDescription(t *testing.T) {
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(testPageHTML))
	assert.Nil(t, err)

	description := GetPageDescriptionHTML(doc)
	assert.Equal(t, description, "Some <b>page</b> description")
}
