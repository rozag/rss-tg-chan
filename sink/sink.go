package sink

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"

	"github.com/rozag/rss-tg-chan/retry"
)

// TODO: refactor

// Sink defines an interface for publishing posts from feeds
type Sink struct {
	config *Config
}

// New constructs a new Sink
func New(config *Config) *Sink {
	return &Sink{config}
}

// Send publishes the post
func (s Sink) Send(feeds map[string][]Post) uint {
	client := &http.Client{
		Timeout: 30 * time.Second,
	}

	cnt := uint(0)
	for _, posts := range feeds {
		for _, post := range posts {
			err := send(client, s.config.TgBotToken, s.config.TgChannel, post)
			if err != nil {
				log.Printf("[ERROR] Failed sending post: %v", err)
			} else {
				cnt++
			}
		}
	}

	return cnt
}

func send(client *http.Client, tgBotToken, tgChannel string, post Post) error {
	// Prepare post's text
	text := post.GetPublishableText()
	if text == "" {
		return fmt.Errorf("Text is empty for the post: %v", post)
	}

	// Prepare the body
	type Body struct {
		ChatID    string `json:"chat_id"`
		ParseMode string `json:"parse_mode"`
		Text      string `json:"text"`
	}
	body := Body{tgChannel, "HTML", text}
	bodyBytes, err := json.Marshal(body)
	if err != nil {
		return err
	}

	// Send the text
	url := fmt.Sprintf("https://api.telegram.org/bot%s/sendMessage", tgBotToken)
	err = retry.Do(3, time.Second, 2, func() error {
		resp, err := client.Post(url, "application/json", bytes.NewReader(bodyBytes))
		if err != nil {
			return err
		}
		defer resp.Body.Close()
		if resp.StatusCode != http.StatusOK {
			bytes, bodyerr := ioutil.ReadAll(resp.Body)
			if bodyerr != nil {
				return fmt.Errorf("Got status code=%d body=%s", resp.StatusCode, string(bodyBytes))
			}
			return fmt.Errorf("Got status code=%d respBody=%s body=%s", resp.StatusCode, string(bytes), string(bodyBytes))
		}
		return nil
	})
	return err
}
