package storage

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/rozag/rss-tg-chan/retry"
)

// Storage knows how to load and save last published time of the feeds
type Storage struct {
	config *Config
}

// LoadTimes maps each feed url to the last published time for the feed
func (s Storage) LoadTimes(urls []string) (map[string]time.Time, error) {
	state, err := loadState(s.config.GithubToken, s.config.GithubGistID, s.config.GithubGistFileName)
	if err != nil {
		return nil, err
	}

	times, err := parseState(state)
	if err != nil {
		log.Printf("[ERROR] Failed to state (%s): %v", state, err)
		return buildZeroedTimes(urls), nil
	}

	return times, nil
}

func loadState(githubToken, githubGistID, githubGistFileName string) (string, error) {
	// Try to load our storage Gist
	client := &http.Client{
		Timeout: 15 * time.Second,
	}
	url := fmt.Sprintf("https://api.github.com/gists/%s", githubGistID)
	req, err := http.NewRequest(http.MethodGet, url, nil)
	authHeader := fmt.Sprintf("token %s", githubToken)
	req.Header.Add("Authorization", authHeader)
	var resp *http.Response
	err = retry.Do(3, time.Second, 2, func() error {
		resp, err = client.Do(req)
		if err != nil {
			return err
		}
		if resp.StatusCode != http.StatusOK {
			resp.Body.Close()
			return fmt.Errorf("Cannot load state. Got status code: %d", resp.StatusCode)
		}
		return nil
	})
	if err != nil {
		log.Printf("[ERROR] Storage loading failed: %v", err)
		return "", err
	}
	defer resp.Body.Close()

	// Get the response body bytes
	bytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	// Parse state string
	type File struct {
		Content string `json:"content"`
	}
	type Root struct {
		Files map[string]File `json:"files"`
	}
	var r Root
	err = json.Unmarshal(bytes, &r)
	if err != nil {
		return "", err
	}

	state := r.Files[githubGistFileName].Content
	return state, nil
}

func parseState(state string) (map[string]time.Time, error) {
	bytes := []byte(state)
	timesUnparsed := make(map[string]string)
	err := json.Unmarshal(bytes, &timesUnparsed)
	if err != nil {
		return nil, err
	}

	times := make(map[string]time.Time, len(timesUnparsed))
	for url, timeStr := range timesUnparsed {
		time, err := time.Parse(time.RFC3339, timeStr)
		if err != nil {
			log.Printf("[ERROR] Failed to parse %s: %v", timeStr, err)
			continue
		}
		times[url] = time
	}

	return times, nil
}

func buildZeroedTimes(urls []string) map[string]time.Time {
	times := make(map[string]time.Time, len(urls))
	for _, url := range urls {
		times[url] = time.Unix(0, 0).UTC()
	}
	return times
}

// SaveTimes saves last published time for each feed
func (s Storage) SaveTimes(times map[string]time.Time) error {
	// Build a map[string]string from the times
	timesStrs := make(map[string]string, len(times))
	for url, time := range times {
		bytes, err := time.MarshalText()
		if err != nil {
			log.Printf("[ERROR] Failed to marshal time (%v): %v", time, err)
			continue
		}
		timesStrs[url] = string(bytes)
	}

	// Get state JSON string from the stringified times
	bytes, err := json.Marshal(timesStrs)
	if err != nil {
		return err
	}
	state := string(bytes)

	// Prepare the body
	bytes, err = json.Marshal(state)
	if err != nil {
		return err
	}
	bodyTemplate := `{"files":{"%s":{"content":%s}}}`
	body := fmt.Sprintf(bodyTemplate, s.config.GithubGistFileName, string(bytes))
	bodyReader := strings.NewReader(body)

	// Save our state
	client := &http.Client{
		Timeout: 15 * time.Second,
	}
	url := fmt.Sprintf("https://api.github.com/gists/%s", s.config.GithubGistID)
	req, err := http.NewRequest(http.MethodPatch, url, bodyReader)
	authHeader := fmt.Sprintf("token %s", s.config.GithubToken)
	req.Header.Add("Authorization", authHeader)
	err = retry.Do(3, time.Second, 2, func() error {
		resp, err := client.Do(req)
		if err != nil {
			return err
		}
		defer resp.Body.Close()
		if resp.StatusCode != http.StatusOK {
			return fmt.Errorf("Cannot save state. Got status code: %d", resp.StatusCode)
		}
		return nil
	})
	return err
}

// New constructs a new Storage
func New(config *Config) *Storage {
	return &Storage{config: config}
}
