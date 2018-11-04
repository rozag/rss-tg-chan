package source

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"

	"github.com/rozag/rss-tg-chan/retry"
)

// TODO: refactor

// Source knows how to load feeds' urls
type Source struct {
	config *Config
	cache  []string
}

// LoadFeeds loads a slice of feeds' urls from the specified URL
func (s Source) LoadFeeds() ([]string, error) {
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
	client := &http.Client{
		Timeout: 30 * time.Second,
	}
	var resp *http.Response
	var err error
	err = retry.Do(3, time.Second, 2, func() error {
		resp, err = client.Get(URL)
		if err != nil {
			return err
		}
		if resp.StatusCode != http.StatusOK {
			resp.Body.Close()
			return fmt.Errorf("Cannot load feeds. Got status code: %d", resp.StatusCode)
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	// Get the response body bytes
	bytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	// Parse Feeds struct from the json
	type Feeds struct {
		Urls []string `json:"feeds"`
	}
	var f Feeds
	err = json.Unmarshal(bytes, &f)
	if err != nil {
		return nil, err
	}

	return f.Urls, nil
}

// New constructs a new Source
func New(config *Config) *Source {
	return &Source{config: config}
}
