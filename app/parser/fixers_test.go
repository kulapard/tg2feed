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
	for _, tb := range tbl {
		doc, err := goquery.NewDocumentFromReader(strings.NewReader(tb.html))
		assert.Nil(t, err)

		body := doc.Find("body")

		fixedBody := FixEmoji(body)

		assert.NotEqualf(t, body, fixedBody, "body and fixedBody should not be equal")

		fixedBodyHTML, err := fixedBody.Html()
		assert.Nil(t, err)

		assert.Equal(t, tb.fixedBody, fixedBodyHTML)
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

func TestRemoveUnsafeTags(t *testing.T) {
	const html = `<body><i>t</i><img src="img.jpg"><b>e</b><a href="link">s<to-remove>t</to-remove></a><br><br/></body>`
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(html))
	assert.Nil(t, err)

	body := doc.Find("body")

	fixedBody := RemoveUnsafeTags(body)

	assert.NotEqualf(t, body, fixedBody, "body and fixedBody should not be equal")

	fixedBodyHTML, err := fixedBody.Html()
	assert.Nil(t, err)

	assert.Equal(t, `<i>t</i><b>e</b><a href="link">st</a><br/><br/>`, fixedBodyHTML)
}
