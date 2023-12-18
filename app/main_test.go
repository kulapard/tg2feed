package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetConfig(t *testing.T) {
	cfg := getConfig()
	assert.Equal(t, cfg.OutputDir, "./")
	assert.Equal(t, cfg.TelegramChannels, []string{"@telegram"})
	assert.Equal(t, cfg.Formats, []string{"rss"})
}
