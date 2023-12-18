// Package feed provides the feed generator for the telegram channel page.
package feed

import (
	"crypto/sha256"
	"encoding/hex"
	"github.com/gorilla/feeds"
	"github.com/kulapard/tg2feed/app/parser"
	"log"
	"os"
	"strings"
	"time"
)

// Merge feeds
func Merge(fs []*feeds.Feed) *feeds.Feed {
	// Extract channels links
	links := make([]string, len(fs))
	for i, feed := range fs {
		links[i] = feed.Link.Href
	}
	linksStr := strings.Join(links, ", ")
	mergedFeed := &feeds.Feed{
		Title:       "Telegram Feed",
		Description: "Channels: " + linksStr,
		Link:        &feeds.Link{Href: "https://github.com/kulapard/tg2feed"},
		Created:     time.Now(),
	}
	// Merge items
	for _, feed := range fs {
		mergedFeed.Items = append(mergedFeed.Items, feed.Items...)
	}

	// Sort items by created date
	sorFunc := func(a, b *feeds.Item) bool {
		return a.Created.After(b.Created)
	}
	mergedFeed.Sort(sorFunc)

	return mergedFeed
}

// GetFeed returns RSS feed for Telegram channel web page
func GetFeed(page *parser.Page) *feeds.Feed {
	now := time.Now()
	feed := &feeds.Feed{
		Title:       page.Title,
		Link:        &feeds.Link{Href: page.Link},
		Description: page.Description,
		Created:     now,
		Updated:     now,
		Author:      &feeds.Author{Name: page.Title},
	}

	if page.ImageURL != "" {
		feed.Image = &feeds.Image{
			Url: page.ImageURL,
		}
	}

	feed.Items = make([]*feeds.Item, len(page.Posts))

	for i, post := range page.Posts {
		var enclosure *feeds.Enclosure
		if len(post.Images) > 0 {
			enclosure = &feeds.Enclosure{
				Url:    post.Images[0],
				Length: "0", //todo: get length
				Type:   "image/jpeg",
			}
		} else if len(post.Videos) > 0 {
			enclosure = &feeds.Enclosure{
				Url:    post.Videos[0],
				Length: "0", //todo: get length
				Type:   "video/mp4",
			}
		}
		feed.Items[i] = &feeds.Item{
			Id:          GetGUID(post.Link),
			Title:       post.Title,
			Link:        &feeds.Link{Href: post.Link},
			Description: post.Text,
			Author:      &feeds.Author{Name: page.Title},
			Created:     post.Created,
		}
		if enclosure != nil {
			feed.Items[i].Enclosure = enclosure
		}
	}

	// Sort items by created date
	sorFunc := func(a, b *feeds.Item) bool {
		return a.Created.After(b.Created)
	}
	feed.Sort(sorFunc)

	return feed
}

// GetGUID returns the GUID for the specified string
func GetGUID(str string) string {
	hash := sha256.Sum256([]byte(str))
	hashStr := hex.EncodeToString(hash[:])
	return hashStr
}

func save(fileName, content string) error {
	fh, err := os.Create(fileName) //nolint:gosec // tolerable security risk
	if err != nil {
		return err
	}
	defer fh.Close() // nolint
	if _, err = fh.WriteString(content); err != nil {
		return err
	}
	log.Printf("[INFO] feed file saved to %s", fileName)
	return nil
}

func saveToRSS(feed *feeds.Feed, dir string) error {
	fname := dir + "/rss.xml"
	content, err := feed.ToRss()
	if err != nil {
		return err
	}
	return save(fname, content)
}

func saveToAtom(feed *feeds.Feed, dir string) error {
	fname := dir + "/atom.xml"
	content, err := feed.ToAtom()
	if err != nil {
		return err
	}
	return save(fname, content)
}

func saveToJSON(feed *feeds.Feed, dir string) error {
	fname := dir + "/feed.json"
	content, err := feed.ToJSON()
	if err != nil {
		return err
	}
	return save(fname, content)
}

// SaveToFile saves RSS feed to file
func SaveToFile(f *feeds.Feed, dir string, formats []string) error {
	// Check id directory exists
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		// Create directory recursively
		err = os.MkdirAll(dir, 0o755) //nolint:gosec // tolerable security risk
		if err != nil {
			return err
		}

		log.Printf("[INFO] directory created: %s", dir)
	}

	// Generate feed string for each format
	for _, format := range formats {
		switch format {
		case "rss":
			err := saveToRSS(f, dir)
			if err != nil {
				return err
			}
		case "atom":
			err := saveToAtom(f, dir)
			if err != nil {
				return err
			}
		case "json":
			err := saveToJSON(f, dir)
			if err != nil {
				return err
			}
		default:
			log.Printf("[ERROR] ignoring unknown format: %s", format)
		}
	}
	return nil
}
