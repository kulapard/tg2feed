package parser

import (
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"github.com/stretchr/testify/assert"
	"strings"
	"testing"
	"time"
)

func TestGetSafeHTML(t *testing.T) {
	const html = `<body>
	<a href="https://t.me/s/telegram" target="_blank" rel="noopener" onclick="return confirm('Open this link?\n\n'+this.href);">telegram</a>
	<i class="emoji">👍</i>
	<tg-emoji>👍</tg-emoji>
</body>`
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(html))
	assert.Nil(t, err)

	body := doc.Find("body")

	safeHTML := GetSafeHTML(body)
	assert.Equal(t, `<a href="https://t.me/s/telegram">telegram</a> 👍 👍`, safeHTML)

	// Empty body
	doc, err = goquery.NewDocumentFromReader(strings.NewReader(""))
	assert.Nil(t, err)
	body = doc.Find("body")
	safeHTML = GetSafeHTML(body)
	assert.Equal(t, "", safeHTML)
}

func TestParseDateTime(t *testing.T) {
	tbl := []struct {
		inp string
		err error
		out string
	}{
		{"", fmt.Errorf("can't parse empty datetime"), time.Now().Format(time.RFC822Z)},
		{"2023-12-13T13:16:01+00:00", nil, "13 Dec 23 13:16 +0000"},
		{"2023-12-13T13:16:01+01:00", nil, "13 Dec 23 13:16 +0100"},
		{"12345", fmt.Errorf("can't parse datetime 12345"), time.Now().Format(time.RFC822Z)},
	}

	for _, tb := range tbl {
		ts, err := ParseDateTime(tb.inp)
		assert.Equal(t, tb.err, err)
		assert.Equal(t, tb.out, ts.Format(time.RFC822Z))
	}
}

func TestShortenText(t *testing.T) {
	tbl := []struct {
		inp string
		max int
		out string
	}{
		{"", 1, ""},
		{"abcdef", 1, "a..."},
		{"abcdef", 3, "abc..."},
		{"abc def", 5, "abc..."},
		{"abc def", 7, "abc def"},
	}

	for _, tb := range tbl {
		short := ShortenText(tb.inp, tb.max)
		assert.Equal(t, tb.out, short)
	}
}
