package main

import (
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/hashicorp/logutils"
	"github.com/mmcdole/gofeed"
	"github.com/rozag/rss-tg-chan/config"
	"github.com/rozag/rss-tg-chan/feed"
	"github.com/rozag/rss-tg-chan/post"
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
		log.Println("[DEBUG] Starting feeds loading")
		feeds, err := feed.LoadFeeds(config.SourcesURL)
		if err == nil {
			log.Printf("[DEBUG] Successfully loaded %d feeds", len(feeds))
			runFeedsProcessing(feeds, config.Workers)
		} else {
			log.Printf("[ERROR] Cannot load feeds: %v", err)
		}
		time.Sleep(config.Period)
	}
}

func runFeedsProcessing(feeds []feed.Feed, workers uint) {
	numJobs := len(feeds)
	jobs := make(chan feed.Feed, numJobs)
	results := make(chan []*post.Post, numJobs)

	for i := uint(0); i < workers; i++ {
		go feedProcessor(jobs, results)
	}

	for _, feed := range feeds {
		jobs <- feed
	}
	close(jobs)

	var posts []*post.Post
	for i := 0; i < numJobs; i++ {
		p := <-results
		posts = append(posts, p...)
	}
	log.Printf("[DEBUG] %d workers successfully loaded %d posts", workers, len(posts))
	// TODO
	for _, post := range posts {
		fmt.Println(*post)
	}
}

func feedProcessor(jobs <-chan feed.Feed, results chan<- []*post.Post) {
	fp := gofeed.NewParser()
	for f := range jobs {
		feed, err := fp.ParseURL(f.URL)
		if err != nil || feed == nil {
			results <- nil
			continue
		}

		var posts []*post.Post
		for _, item := range feed.Items {
			if item == nil {
				continue
			}
			post := post.New(
				strings.TrimSpace(item.Title),
				strings.TrimSpace(item.Description),
				strings.TrimSpace(item.Link),
			)
			posts = append(posts, post)
		}
		results <- posts
	}
}
