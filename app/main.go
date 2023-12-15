package main

import (
	"fmt"
	"github.com/gorilla/feeds"
	"github.com/kulapard/tg2rss/app/feed"
	"github.com/kulapard/tg2rss/app/parser"
	"log"
	"os"
	"strings"
)

var revision = "unknown"

const defaultOutputFile = "rss.xml"

// Config represents application configuration
type Config struct {
	OutputFile       string
	TelegramChannels []string
}

func (c *Config) String() string {
	return fmt.Sprintf("OutputFile: %s\nTelegramChannels: %v", c.OutputFile, c.TelegramChannels)
}

func getConfig() *Config {
	// Set default output file
	outfile := os.Getenv("INPUT_OUTPUT-FILE")
	if outfile == "" {
		outfile = defaultOutputFile
	}
	return &Config{
		OutputFile:       outfile,
		TelegramChannels: strings.Split(os.Getenv("INPUT_TELEGRAM-CHANNELS"), ","),
	}
}

func main() {
	fmt.Println("Running tg2rss " + revision)
	cfg := getConfig()

	// Print config
	Info(cfg.String())

	var tgFeed *feeds.Feed
	tgFeeds := make([]*feeds.Feed, len(cfg.TelegramChannels))

	// Build RSS feed for each channel
	for i, tgChannel := range cfg.TelegramChannels {
		Info("Building RSS feed for Telegram channel: " + tgChannel)
		// Parse the page
		page := parser.Parse(tgChannel)
		tgFeeds[i] = feed.GetFeed(&page)
	}

	// Merge all feeds if there are more than one
	if len(tgFeeds) > 1 {
		// Merge all feeds
		tgFeed = feed.Merge(tgFeeds)
		Info(fmt.Sprintf("Merged %d RSS feeds", len(tgFeeds)))
	} else {
		tgFeed = tgFeeds[0]
	}

	// Save RSS feed to file
	err := feed.SaveToFile(tgFeed, cfg.OutputFile)
	if err != nil {
		log.Fatal(err)
	}
	Info("RSS feed saved to file: " + cfg.OutputFile)
}
