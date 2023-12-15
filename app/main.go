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
	DryRun           bool
	OutputFile       string
	TelegramChannels []string
}

func (c *Config) String() string {
	return fmt.Sprintf("DryRun: %v\nOutputFile: %s\nTelegramChannels: %v", c.DryRun, c.OutputFile, c.TelegramChannels)
}

// str2Bool converts string to bool
func str2Bool(str string) bool {
	return str == "true"
}

func getConfig() *Config {
	// Set default output file
	outfile := os.Getenv("INPUT_OUTPUT-FILE")
	if outfile == "" {
		outfile = defaultOutputFile
	}
	return &Config{
		DryRun:           str2Bool(os.Getenv("INPUT_DRY-RUN")),
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
		Info("Telegram channel: " + tgChannel)
		// Parse the page
		page := parser.Parse(tgChannel)
		tgFeeds[i] = feed.GetFeed(&page)
	}

	if len(tgFeeds) > 1 {
		// Merge all feeds
		tgFeed = feed.Merge(tgFeeds)
	} else {
		tgFeed = tgFeeds[0]
	}

	// Save RSS feed to file
	if !cfg.DryRun {
		err := feed.SaveToFile(tgFeed, cfg.OutputFile)
		if err != nil {
			log.Fatal(err)
		}
	}
}
