package source

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"time"

	"github.com/rozag/rss-tg-chan/retry"
	"github.com/rozag/rss-tg-chan/timeout"
)

// Source is an interface for feeds' urls provider
type Source interface {
	LoadFeeds() ([]string, error)
}

type source struct {
	config *Config
	cache  []string
}

// Feeds contains a slice of feeds' urls
type Feeds struct {
	Urls []string `json:"feeds"`
}

// LoadFeeds loads a slice of feeds' urls from the specified URL
func (s source) LoadFeeds() ([]string, error) {
	// Try to load feeds
	urls, err := loadFeeds(s.config.SourcesURL)
	if err != nil {
		// Return error if cache is empty
		if s.cache == nil {
			return nil, err
		}
		// Return cache if there is some
		log.Printf("[ERROR] Failed loading feeds, returning cache. %v", err)
		return s.cache, nil
	}

	// Put the loaded feeds into the cache
	s.cache = urls

	return urls, nil
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

// New constructs a new Source
func New(config *Config) Source {
	return source{config: config}
}
