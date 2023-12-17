package parser

import (
	"github.com/PuerkitoBio/goquery"
	"github.com/stretchr/testify/assert"
	"strings"
	"testing"
)

func TestFixEmoji(t *testing.T) {
	tbl := []struct {
		html      string
		fixedBody string
	}{
		{`<html><body><p><i class="emoji">👍</i></p></body></html>`, `<p>👍</p>`},
		{`<html><body><p><tg-emoji>👍</tg-emoji></p></body></html>`, `<p>👍</p>`},
	}
	for _, tt := range tbl {
		doc, err := goquery.NewDocumentFromReader(strings.NewReader(tt.html))
		assert.Nil(t, err)

		body := doc.Find("body")

		fixedBody := FixEmoji(body)

		assert.NotEqualf(t, body, fixedBody, "body and fixedBody should not be equal")

		fixedBodyHTML, err := fixedBody.Html()
		assert.Nil(t, err)

		assert.Equal(t, tt.fixedBody, fixedBodyHTML)
	}
}

func TestFixLinks(t *testing.T) {
	const html = `<body><a href="https://t.me/s/telegram" target="_blank" rel="noopener" onclick="return confirm('Open this link?\n\n'+this.href);">telegram</a></body>`
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(html))
	assert.Nil(t, err)

	body := doc.Find("body")

	fixedBody := FixLinks(body)

	assert.NotEqualf(t, body, fixedBody, "body and fixedBody should not be equal")

	fixedBodyHTML, err := fixedBody.Html()
	assert.Nil(t, err)

	assert.Equal(t, `<a href="https://t.me/s/telegram">telegram</a>`, fixedBodyHTML)
}
