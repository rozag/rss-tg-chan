package main

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/hashicorp/logutils"
	"github.com/rozag/rss-tg-chan/config"
	"github.com/rozag/rss-tg-chan/feed"
)

func main() {
	// Load the flags into the config.Config struct
	config, err := config.ParseFlags()
	if err != nil {
		log.Printf("[ERROR] Cannot load config: %v", err)
		return
	}

	// Set up the logs filter
	filter := &logutils.LevelFilter{
		Levels:   []logutils.LogLevel{"DEBUG", "ERROR"},
		MinLevel: logutils.LogLevel(config.LogLevel),
		Writer:   os.Stderr,
	}
	log.SetOutput(filter)

	// Print debug config info
	log.Printf("[DEBUG] Using log level: %s", config.LogLevel)
	log.Printf("[DEBUG] Using sources url: %s", config.SourcesURL)
	log.Printf("[DEBUG] Using period: %v", config.Period)

	// Run the fetch-and-post loop
	done := make(chan bool)
	go runFeedsLoop(done, config)
	<-done
}

func runFeedsLoop(done chan<- bool, config *config.Config) {
	for {
		feeds, err := feed.LoadFeeds(config.SourcesURL)
		if err == nil {
			runFeedsProcessing(feeds)
		} else {
			log.Printf("[ERROR] Cannot load feeds: %v", err)
		}
		time.Sleep(config.Period)
	}
}

func runFeedsProcessing(feeds []feed.Feed) {
	for _, feed := range feeds {
		fmt.Println(feed)
	}
}
