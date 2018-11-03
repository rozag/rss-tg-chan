package app

import (
	"log"
	"time"

	"github.com/mmcdole/gofeed"
	"github.com/rozag/rss-tg-chan/sink"
	"github.com/rozag/rss-tg-chan/source"
	"github.com/rozag/rss-tg-chan/storage"
)

// App contains all parts of the service
type App struct {
	config  *Config
	source  *source.Source
	storage *storage.Storage
	sink    *sink.Sink
}

// New constructs a new App instance
func New(config *Config) *App {
	return &App{
		config:  config,
		source:  source.New(config.SourceConfig),
		storage: storage.New(config.StorageConfig),
		sink:    sink.New(config.SinkConfig),
	}
}

// Run starts the main app loop
func (app *App) Run(done chan<- bool) {
	for {
		app.run()
		if app.config.SingleRun {
			break
		}
		time.Sleep(app.config.Period)
	}
	done <- true
}

func (app *App) run() {
	// Load feeds' urls
	log.Println("[DEBUG] Starting feeds loading")
	urls, err := app.source.LoadFeeds()
	if err != nil {
		log.Printf("[ERROR] Cannot load feeds: %v", err)
		return
	}
	log.Printf("[DEBUG] Successfully loaded %d feeds", len(urls))

	// Load app state
	times, err := app.storage.LoadTimes(urls)
	if err != nil {
		log.Printf("[ERROR] Failed to load state: %v", err)
		return
	}

	// Load posts from every feed
	posts, postsCnt := loadPosts(urls, app.config.Workers)
	log.Printf("[DEBUG] %d workers successfully loaded %d posts", app.config.Workers, postsCnt)

	// Filter out outdated posts
	filtered := make(map[string][]sink.Post, len(posts))
	filteredCnt := uint(0)
	for url, ps := range posts {
		if lastPublished, ok := times[url]; ok {
			for _, post := range ps {
				if post.PublishedAfter(lastPublished) {
					filtered[url] = append(filtered[url], post)
					filteredCnt++
				}
			}
		} else {
			filtered[url] = append(filtered[url], ps...)
			filteredCnt += uint(len(ps))
		}
	}
	log.Printf("[DEBUG] Posts count after filtering out outdated: %d", filteredCnt)
}

type batch struct {
	url   string
	posts []sink.Post
}

func loadPosts(urls []string, workers uint) (map[string][]sink.Post, uint) {
	numJobs := len(urls)
	jobs := make(chan string, numJobs)
	results := make(chan batch, numJobs)

	for i := uint(0); i < workers; i++ {
		go feedLoader(jobs, results)
	}

	for _, url := range urls {
		jobs <- url
	}
	close(jobs)

	var batches []batch
	for i := 0; i < numJobs; i++ {
		b := <-results
		batches = append(batches, b)
	}

	count := uint(0)
	posts := make(map[string][]sink.Post, len(batches))
	for _, batch := range batches {
		count += uint(len(batch.posts))
		posts[batch.url] = append(posts[batch.url], batch.posts...)
	}

	return posts, count
}

func feedLoader(jobs <-chan string, results chan<- batch) {
	parser := gofeed.NewParser()
	for url := range jobs {
		feed, err := parser.ParseURL(url)
		if err != nil {
			log.Printf("[ERROR] Failed to load feed from url=%s: %v", url, err)
			continue
		}
		if feed == nil {
			log.Printf("[ERROR] Loaded feed is nil; url=%s", url)
			continue
		}

		var posts []sink.Post
		for _, item := range feed.Items {
			if item == nil {
				continue
			}
			post := sink.NewPost(
				item.Title,
				item.Description,
				item.Link,
				item.PublishedParsed,
			)
			posts = append(posts, *post)
		}
		results <- batch{url, posts}
	}
}
