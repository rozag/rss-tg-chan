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
	source  source.Source
	storage storage.Storage
	sink    sink.Sink
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
		time.Sleep(app.config.Period)
	}
}

type publishedFeed struct {
	url           string
	lastPublished time.Time
}

func (app *App) run() {
	log.Println("[DEBUG] Starting feeds loading")
	urls, err := app.source.LoadFeeds()
	if err != nil {
		log.Printf("[ERROR] Cannot load feeds: %v", err)
		return
	}
	log.Printf("[DEBUG] Successfully loaded %d feeds", len(urls))

	app.storage.Update()
	feeds := make([]publishedFeed, len(urls))
	for i, url := range urls {
		feeds[i] = publishedFeed{url: url, lastPublished: app.storage.GetLastPublishedTime(url)}
	}

	runFeedsProcessing(feeds, app.config.Workers)
}

func runFeedsProcessing(feeds []publishedFeed, workers uint) {
	numJobs := len(feeds)
	jobs := make(chan publishedFeed, numJobs)
	results := make(chan []*sink.Post, numJobs)

	for i := uint(0); i < workers; i++ {
		go feedProcessor(jobs, results)
	}

	for _, pf := range feeds {
		jobs <- pf
	}
	close(jobs)

	var posts []*sink.Post
	for i := 0; i < numJobs; i++ {
		p := <-results
		posts = append(posts, p...)
	}
	log.Printf("[DEBUG] %d workers successfully loaded %d posts", workers, len(posts))
	// TODO
	// for _, post := range posts {
	// 	fmt.Println(*post)
	// }
}

func feedProcessor(jobs <-chan publishedFeed, results chan<- []*sink.Post) {
	fp := gofeed.NewParser()
	for pf := range jobs {
		feed, err := fp.ParseURL(pf.url)
		if err != nil || feed == nil {
			results <- nil
			continue
		}

		// TODO: filter out posts with published time before last published

		var posts []*sink.Post
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
			posts = append(posts, post)
		}
		results <- posts
	}
}
