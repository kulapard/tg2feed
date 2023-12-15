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

// Str2Bool converts string to bool
func Str2Bool(str string) bool {
	return str == "true"
}

func main() {
	fmt.Println("Running tg2rss " + revision)

	// Print all environment variables
	for _, e := range os.Environ() {
		pair := strings.Split(e, "=")
		Info(pair[0] + "=" + pair[1])
	}

	dryRun := Str2Bool(os.Getenv("INPUT_DRY-RUN"))
	outputFile := os.Getenv("INPUT_OUTPUT-FILE")
	tgChannelsStr := os.Getenv("INPUT_TELEGRAM-CHANNELS")

	if dryRun {
		Info("Dry run")
	}
	Info("Output file: " + outputFile)
	Info("Telegram channel: " + tgChannelsStr)

	// Split channel name by comma
	tgChannels := strings.Split(tgChannelsStr, ",")

	var tgFeed *feeds.Feed
	tgFeeds := make([]*feeds.Feed, len(tgChannels))

	// Build RSS feed for each channel
	for i, tgChannel := range tgChannels {
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
	// Set default output file
	if outputFile == "" {
		Warning("Output file is not set. Using default: " + defaultOutputFile)
		outputFile = defaultOutputFile
	}

	// Save RSS feed to file
	if !dryRun {
		err := feed.SaveToFile(tgFeed, outputFile)
		if err != nil {
			log.Fatal(err)
		}
	}
}
