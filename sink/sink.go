package sink

import "fmt"

func (post Post) String() string {
	return fmt.Sprintf(
		"Post{\n\tTitle=%s\n\tDescription=%s\n\tURL=%s\n\tPublished=%v\n}",
		post.Title,
		post.Description,
		post.URL,
		post.Published,
	)
}

// Sink defines an interface for publishing posts from feeds
type Sink interface {
}

// New constructs a new Sink
func New(config *Config) Sink {
	return nil
}
