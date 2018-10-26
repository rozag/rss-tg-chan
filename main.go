package main

import (
	"fmt"
	"log"
	"os"

	"github.com/hashicorp/logutils"
	"github.com/rozag/rss-tg-chan/config"
	"github.com/rozag/rss-tg-chan/feed"
)

func main() {
	config, err := config.ParseFlags()
	if err != nil {
		log.Printf("[ERROR] Cannot load config: %v", err)
	}

	filter := &logutils.LevelFilter{
		Levels:   []logutils.LogLevel{"DEBUG", "ERROR"},
		MinLevel: logutils.LogLevel(config.LogLevel),
		Writer:   os.Stderr,
	}
	log.SetOutput(filter)
	log.Printf("[DEBUG] Using log level: %s", config.LogLevel)

	log.Printf("[DEBUG] Using sources url: %s", config.SourcesURL)

	// TODO
	feeds, err := feed.LoadFeeds(config.SourcesURL)
	if err != nil {
		log.Printf("[ERROR] Cannot load feeds: %v", err)
		return
	}
	for _, feed := range feeds {
		fmt.Println(feed)
	}
}
