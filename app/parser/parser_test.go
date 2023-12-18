package parser

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestGetChannelWebURL(t *testing.T) {
	tbl := []struct {
		chName string
		chURL  string
	}{
		{`telegram`, `https://t.me/s/telegram`},
		{`@telegram`, `https://t.me/s/telegram`},
		{`https://t.me/telegram`, `https://t.me/s/telegram`},
		{`https://t.me/s/telegram`, `https://t.me/s/telegram`},
		{`https://t.me/s/telegram?foo=bar#extra`, `https://t.me/s/telegram`},
	}
	for _, tb := range tbl {
		url := GetChannelWebURL(tb.chName)
		assert.Equal(t, tb.chURL, url)
	}
}
