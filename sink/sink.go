package sink

import (
	"fmt"
	"log"
	"net/http"
	"net/url"
	"time"

	"github.com/rozag/rss-tg-chan/retry"
)

// Sink defines an interface for publishing posts from feeds
type Sink struct {
	config *Config
}

// New constructs a new Sink
func New(config *Config) *Sink {
	return &Sink{config}
}

// Send publishes the post
func (s Sink) Send(feeds map[string][]Post) {
	client := &http.Client{
		Timeout: 15 * time.Second,
	}
	for _, posts := range feeds {
		for _, post := range posts {
			send(client, s.config.TgBotToken, s.config.TgChannel, post)
		}
	}
}

func send(client *http.Client, tgBotToken, tgChannel string, post Post) {
	// Prepare post's text
	text := post.GetPublishableText()
	if text == "" {
		return
	}

	// Prepare URL
	link, err := url.Parse("https://api.telegram.org")
	if err != nil {
		log.Printf("[ERROR] Cannot parse URL. Got error: %v", err)
		return
	}
	link.Path += fmt.Sprintf("/bot%s/sendMessage", tgBotToken)
	params := url.Values{}
	params.Add("chat_id", tgChannel)
	params.Add("parse_mode", "Markdown")
	params.Add("text", text)
	link.RawQuery = params.Encode()

	// Send the text
	err = retry.Do(3, time.Second, 2, func() error {
		resp, err := client.Post(link.String(), "", nil)
		if err != nil {
			return err
		}
		defer resp.Body.Close()
		if resp.StatusCode != http.StatusOK {
			return fmt.Errorf("Got status code: %d", resp.StatusCode)
		}
		return nil
	})
	if err != nil {
		log.Printf("[ERROR] Cannot send the post. Got error: %v", err)
		return
	}
}
