package feed

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
)

// Feed contains RSS feed info
type Feed struct {
	Title string `json:"title"`
	URL   string `json:"url"`
}

// Feeds contains a slice of Feed structs
type Feeds struct {
	Feeds []Feed `json:"feeds"`
}

// LoadFeeds loads a slice of Feed structs from the specified URL
func LoadFeeds(URL string) ([]Feed, error) {
	resp, err := http.Get(URL)
	if err != nil {
		return nil, err
	}

	bytes, err := ioutil.ReadAll(resp.Body)

	var f Feeds
	err = json.Unmarshal(bytes, &f)
	if err != nil {
		return nil, err
	}

	return f.Feeds, nil
}
