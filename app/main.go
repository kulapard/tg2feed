package main

import (
	"fmt"
	"github.com/gorilla/feeds"
	"github.com/kulapard/tg2feed/app/feed"
	"github.com/kulapard/tg2feed/app/parser"
	"log"
	"os"
	"strings"
)

var revision = "unknown"

const defaultOutputDir = "./"
const defaultFormat = "rss"
const defaultChannel = "@telegram"

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
	// Set default output dir
	outdir := os.Getenv("INPUT_OUTPUT-DIR")
	if outdir == "" {
		outdir = defaultOutputDir
	}
	// Set default Telegram channels
	chStr := os.Getenv("INPUT_TELEGRAM-CHANNELS")
	if chStr == "" {
		chStr = defaultChannel
	}
	channels := strings.Split(chStr, ",")

	// Set default formats
	formatStr := os.Getenv("INPUT_FORMATS")
	if formatStr == "" {
		formatStr = defaultFormat
	}

	formats := strings.Split(formatStr, ",")

	return &Config{
		OutputDir:        outdir,
		TelegramChannels: channels,
		Formats:          formats,
	}
}

func main() {
	fmt.Println("Running tg2feed " + revision)
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

	if tgFeed == nil {
		log.Fatal("RSS feed is empty")
	}

	// Save RSS feed to file
	err := feed.SaveToFile(tgFeed, cfg.OutputDir, cfg.Formats)
	if err != nil {
		log.Fatal(err)
	}
}
