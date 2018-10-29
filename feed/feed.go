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
var cache []string

// Feeds contains a slice of Feed structs
type Feeds struct {
	Urls []string `json:"feeds"`
}

func loadFeeds(URL string) ([]string, error) {
	// Try to load data with retry
	var resp *http.Response
	var err error
	err = retry.Do(3, time.Second, 2, func() error {
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

	return f.Urls, nil
}

// LoadFeeds loads a slice of Feed structs from the specified URL
func LoadFeeds(URL string) ([]string, error) {
	// Try to load feeds
	urls, err := loadFeeds(URL)
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
	cache = urls

	return urls, nil
}
