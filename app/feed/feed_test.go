// Package feed provides the feed generator for the telegram channel page.
package feed

import (
	"github.com/gorilla/feeds"
	"github.com/kulapard/tg2feed/app/parser"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestGetFeed(t *testing.T) {
	now := time.Now()
	page := &parser.Page{
		Title:       "Channel Title",
		Link:        "https://t.me/s/telegram",
		Description: "Telegram channel description",
		ImageURL:    "https://telegram.org/img/t_logo.png",
		Posts: []*parser.Post{
			{Title: "Post 2", Link: "https://t.me/s/telegram/2", Text: "Post 2 text", Created: now.Add(time.Hour * -2), Images: []string{"https://telegram.org/img/2.png"}},
			{Title: "Post 1", Link: "https://t.me/s/telegram/1", Text: "Post 1 text", Created: now.Add(time.Hour * -1)},
		},
	}
	feed := GetFeed(page)
	assert.NotNil(t, feed)
	assert.Equal(t, "Channel Title", feed.Title)
	assert.Equal(t, "Channel Title", feed.Author.Name)
	assert.Equal(t, "https://t.me/s/telegram", feed.Link.Href)
	assert.Equal(t, "Telegram channel description", feed.Description)
	assert.Equal(t, "https://telegram.org/img/t_logo.png", feed.Image.Url)

	assert.Equal(t, 2, len(feed.Items))

	item1 := feed.Items[0]
	assert.Equal(t, "Post 1", item1.Title)
	assert.Equal(t, "Post 1 text", item1.Description)
	assert.Equal(t, "https://t.me/s/telegram/1", item1.Link.Href)
	assert.Nil(t, item1.Enclosure)

	item2 := feed.Items[1]
	assert.Equal(t, "Post 2", item2.Title)
	assert.Equal(t, "Post 2 text", item2.Description)
	assert.Equal(t, "https://t.me/s/telegram/2", item2.Link.Href)
	assert.Equal(t, "https://telegram.org/img/2.png", item2.Enclosure.Url)
}

func TestGetGUID(t *testing.T) {
	tbl := []struct {
		inp string
		out string
	}{
		{"", "e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855"},
		{"12345", "5994471abb01112afcc18159f6cc74b4f511b99806da59b3caf5a9c173cacfc5"},
		{"1234567890", "c775e7b757ede630cd0aa1113bd102661ab38829ca52a6422ab782862f268646"},
	}

	for _, tb := range tbl {
		guid := GetGUID(tb.inp)
		assert.Equal(t, tb.out, guid)
		assert.Equal(t, 64, len(guid))
	}
}

func TestMerge(t *testing.T) {
	feed1 := &feeds.Feed{
		Title: "Channel 1",
		Link:  &feeds.Link{Href: "https://t.me/s/telegram1"},
		Items: []*feeds.Item{
			{Title: "Post 1", Link: &feeds.Link{Href: "https://t.me/s/telegram1/1"}},
			{Title: "Post 2", Link: &feeds.Link{Href: "https://t.me/s/telegram1/2"}},
		},
	}
	feed2 := &feeds.Feed{
		Title: "Channel 2",
		Link:  &feeds.Link{Href: "https://t.me/s/telegram2"},
		Items: []*feeds.Item{
			{Title: "Post 3", Link: &feeds.Link{Href: "https://t.me/s/telegram2/3"}},
			{Title: "Post 4", Link: &feeds.Link{Href: "https://t.me/s/telegram2/4"}},
		},
	}
	feed3 := &feeds.Feed{
		Title: "Channel 3",
		Link:  &feeds.Link{Href: "https://t.me/s/telegram3"},
		Items: []*feeds.Item{
			{Title: "Post 5", Link: &feeds.Link{Href: "https://t.me/s/telegram3/5"}},
			{Title: "Post 6", Link: &feeds.Link{Href: "https://t.me/s/telegram3/6"}},
		},
	}
	feed := Merge([]*feeds.Feed{feed1, feed2, feed3})
	assert.NotNil(t, feed)
	assert.Equal(t, "Telegram Feed", feed.Title)
	assert.Equal(t, "Channels: https://t.me/s/telegram1, https://t.me/s/telegram2, https://t.me/s/telegram3", feed.Description)
	assert.Equal(t, 6, len(feed.Items))
}

func TestMerge_Empty(t *testing.T) {
	feed := Merge([]*feeds.Feed{})
	assert.NotNil(t, feed)
	assert.Equal(t, "Telegram Feed", feed.Title)
	assert.Equal(t, "Channels: ", feed.Description)
	assert.Equal(t, 0, len(feed.Items))
}

func TestMerge_Sort(t *testing.T) {
	now := time.Now()
	feed1 := &feeds.Feed{
		Title: "Channel 1",
		Link:  &feeds.Link{Href: "https://t.me/s/telegram1"},
		Items: []*feeds.Item{
			{Title: "Post 2", Link: &feeds.Link{Href: "https://t.me/s/telegram1/2"}, Created: now.Add(time.Hour * -2)},
			{Title: "Post 1", Link: &feeds.Link{Href: "https://t.me/s/telegram1/1"}, Created: now.Add(time.Hour * -1)},
		},
	}
	feed2 := &feeds.Feed{
		Title: "Channel 2",
		Link:  &feeds.Link{Href: "https://t.me/s/telegram2"},
		Items: []*feeds.Item{
			{Title: "Post 4", Link: &feeds.Link{Href: "https://t.me/s/telegram2/4"}, Created: now.Add(time.Hour * -4)},
			{Title: "Post 3", Link: &feeds.Link{Href: "https://t.me/s/telegram2/3"}, Created: now.Add(time.Hour * -3)},
		},
	}
	feed3 := &feeds.Feed{
		Title: "Channel 3",
		Link:  &feeds.Link{Href: "https://t.me/s/telegram3"},
		Items: []*feeds.Item{
			{Title: "Post 6", Link: &feeds.Link{Href: "https://t.me/s/telegram3/6"}, Created: now.Add(time.Hour * -6)},
			{Title: "Post 5", Link: &feeds.Link{Href: "https://t.me/s/telegram3/5"}, Created: now.Add(time.Hour * -4)},
		},
	}
	feed := Merge([]*feeds.Feed{feed1, feed2, feed3})
	assert.NotNil(t, feed)
	assert.Equal(t, "Telegram Feed", feed.Title)
	assert.Equal(t, "Channels: https://t.me/s/telegram1, https://t.me/s/telegram2, https://t.me/s/telegram3", feed.Description)

	// Check items order (from newest to oldest)
	assert.Equal(t, "Post 1", feed.Items[0].Title)
	assert.Equal(t, "Post 2", feed.Items[1].Title)
	assert.Equal(t, "Post 3", feed.Items[2].Title)
	assert.Equal(t, "Post 4", feed.Items[3].Title)
	assert.Equal(t, "Post 5", feed.Items[4].Title)
	assert.Equal(t, "Post 6", feed.Items[5].Title)
}
