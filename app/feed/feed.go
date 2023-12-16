// Package feed provides the feed generator for the telegram channel page.
package feed

import (
	"github.com/gorilla/feeds"
	"github.com/kulapard/tg2rss/app/parser"
	"log"
	"os"
	"sort"
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
		Title:       "Telegram->RSS",
		Description: "Telegram channels: " + linksStr,
		Link:        &feeds.Link{Href: "https://github.com/kulapard/tg2rss"},
		Created:     time.Now(),
	}
	// Merge items
	for _, feed := range fs {
		mergedFeed.Items = append(mergedFeed.Items, feed.Items...)
	}

	// Sort items by date
	sorFunc := func(a, b *feeds.Item) bool {
		return a.Created.After(b.Created)
	}
	mergedFeed.Sort(sorFunc)

	return mergedFeed
}

// GetFeed returns RSS feed for Telegram channel web page
func GetFeed(page *parser.PageData) *feeds.Feed {
	now := time.Now()
	feed := &feeds.Feed{
		Title:       page.Title,
		Link:        &feeds.Link{Href: page.Link},
		Description: page.Description,
		Created:     now,
		Author:      &feeds.Author{Name: page.Title},
	}

	if page.ImageURL != "" {
		feed.Image = &feeds.Image{
			Url: page.ImageURL,
		}
	}

	feed.Items = make([]*feeds.Item, len(page.Posts))

	// Sort posts by date
	sort.Slice(page.Posts, func(i, j int) bool {
		return page.Posts[i].Created.After(page.Posts[j].Created)
	})

	for i, post := range page.Posts {
		var enclosure *feeds.Enclosure
		if len(post.Images) > 0 {
			enclosure = &feeds.Enclosure{
				Url:    post.Images[0].URL,
				Length: "0", //todo: get length
				Type:   "image/jpeg",
			}
		} else if len(post.Previews) > 0 {
			if post.Previews[0].VideoURL != "" {
				enclosure = &feeds.Enclosure{
					Url:    post.Previews[0].VideoURL,
					Length: "0", //todo: get length
					Type:   "video/mp4",
				}
			} else if post.Previews[0].ImageURL != "" {
				enclosure = &feeds.Enclosure{
					Url:    post.Previews[0].ImageURL,
					Length: "0", //todo: get length
					Type:   "image/jpeg",
				}
			}
		}
		if enclosure == nil && post.Video.URL != "" {
			enclosure = &feeds.Enclosure{
				Url:    post.Video.URL,
				Length: "0", //todo: get length
				Type:   "video/mp4",
			}
		}

		feed.Items[i] = &feeds.Item{
			Id:          post.ID,
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
	return feed
}

func save(fname, content string) error {
	fh, err := os.Create(fname) //nolint:gosec // tolerable security risk
	if err != nil {
		return err
	}
	defer fh.Close() // nolint
	if _, err = fh.WriteString(content); err != nil {
		return err
	}
	log.Printf("[INFO] feed file saved to %s", fname)
	return nil
}

func saveToXML(feed *feeds.Feed, dir string) error {
	fname := dir + "/feed.xml"
	content, err := feed.ToRss()
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
		case "xml":
			err := saveToXML(f, dir)
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
