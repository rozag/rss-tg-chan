package feed

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"time"

	"github.com/rozag/rss-tg-chan/retry"
	"github.com/rozag/rss-tg-chan/timeout"
)

// Simple in-memory cache for the feeds
var cache []Feed

// Feed contains RSS feed info
type Feed struct {
	Title string `json:"title"`
	URL   string `json:"url"`
}

// Feeds contains a slice of Feed structs
type Feeds struct {
	Feeds []Feed `json:"feeds"`
}

func loadFeeds(URL string) ([]Feed, error) {
	// Try to load data with retry
	var resp *http.Response
	var err error
	err = retry.Do(&retry.Policy{Retries: 3}, func() error {
		return timeout.Do(15*time.Second, func() error {
			resp, err = http.Get(URL)
			return err
		})
	})
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	// Get the response body bytes
	bytes, err := ioutil.ReadAll(resp.Body)

	// Parse Feeds struct from the json
	var f Feeds
	err = json.Unmarshal(bytes, &f)
	if err != nil {
		return nil, err
	}

	return f.Feeds, nil
}

// LoadFeeds loads a slice of Feed structs from the specified URL
func LoadFeeds(URL string) ([]Feed, error) {
	// Try to load feeds
	feeds, err := loadFeeds(URL)
	if err != nil {
		// Return error if cache is empty
		if cache == nil {
			return nil, err
		}
		// Return cache if there is some
		log.Printf("[ERROR] Failed loading feeds, returning cache. %v", err)
		return cache, nil
	}

	// Put the loaded feeds into the cache
	cache = feeds

	return feeds, nil
}
