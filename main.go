package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/hashicorp/logutils"
	"github.com/mmcdole/gofeed"
	"github.com/rozag/rss-tg-chan/app"
	"github.com/rozag/rss-tg-chan/feed"
	"github.com/rozag/rss-tg-chan/post"
)

func main() {
	config := app.NewConfig()
	config.RegisterFlags()
	flag.Parse()
	config.ValidateFlags()

	// Set up the logs filter
	filter := logutils.LevelFilter{
		Levels:   []logutils.LogLevel{"DEBUG", "ERROR"},
		MinLevel: logutils.LogLevel(config.LogLevel),
		Writer:   os.Stderr,
	}
	log.SetOutput(&filter)

	// Print debug config info
	log.Printf("[DEBUG] %v", config)

	// Run the fetch-and-post loop
	done := make(chan bool)
	go runFeedsLoop(done, config)
	<-done
}

func runFeedsLoop(done chan<- bool, config *app.Config) {
	for {
		log.Println("[DEBUG] Starting feeds loading")
		urls, err := feed.LoadFeeds(config.SourceConfig.SourcesURL)
		if err == nil {
			log.Printf("[DEBUG] Successfully loaded %d feeds", len(urls))
			runFeedsProcessing(urls, config.Workers)
		} else {
			log.Printf("[ERROR] Cannot load feeds: %v", err)
		}
		time.Sleep(config.Period)
	}
}

func runFeedsProcessing(urls []string, workers uint) {
	numJobs := len(urls)
	jobs := make(chan string, numJobs)
	results := make(chan []*post.Post, numJobs)

	for i := uint(0); i < workers; i++ {
		go feedProcessor(jobs, results)
	}

	for _, url := range urls {
		jobs <- url
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

func feedProcessor(jobs <-chan string, results chan<- []*post.Post) {
	fp := gofeed.NewParser()
	for url := range jobs {
		feed, err := fp.ParseURL(url)
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
				item.Title,
				item.Description,
				item.Link,
				item.PublishedParsed,
			)
			posts = append(posts, post)
		}
		results <- posts
	}
}
