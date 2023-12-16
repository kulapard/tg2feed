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

const defaultOutputDir = "./"

// Config represents application configuration
type Config struct {
	OutputDir        string
	TelegramChannels []string
	Formats          []string
}

func (c *Config) String() string {
	return fmt.Sprintf("OutputDir: %s, TelegramChannels: %s, Formats: %s", c.OutputDir, c.TelegramChannels, c.Formats)
}

func getConfig() *Config {
	// Set default output file
	outdir := os.Getenv("INPUT_OUTPUT-DIR")
	if outdir == "" {
		outdir = defaultOutputDir
	}
	return &Config{
		OutputDir:        outdir,
		TelegramChannels: strings.Split(os.Getenv("INPUT_TELEGRAM-CHANNELS"), ","),
		Formats:          strings.Split(os.Getenv("INPUT_FORMATS"), ","),
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
	err := feed.SaveToFile(tgFeed, cfg.OutputDir, cfg.Formats)
	if err != nil {
		log.Fatal(err)
	}
}
