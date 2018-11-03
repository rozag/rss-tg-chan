package sink

import (
	"fmt"
	"html"
	"strings"
	"time"
)

// Post contains a title, a description and a url of a news
type Post struct {
	title       string
	description string
	url         string
	published   *time.Time
}

// PublishedAfter returns true if the post was published after the specified time
func (p Post) PublishedAfter(t time.Time) bool {
	return p.published.After(t)
}

func (p Post) String() string {
	return fmt.Sprintf(
		"Post{\n\ttitle=%s\n\tdescription=%s\n\turl=%s\n\tpublished=%v\n}",
		p.title,
		p.description,
		p.url,
		p.published,
	)
}

func clearString(s string) string {
	nospace := strings.TrimSpace(s)
	unescaped := html.UnescapeString(nospace)
	return unescaped
}

// NewPost returns a pointer to the newly created Post struct
func NewPost(title, description, url string, published *time.Time) *Post {
	if published == nil {
		t := time.Unix(0, 0).UTC()
		published = &t
	}
	return &Post{
		clearString(title),
		clearString(description),
		clearString(url),
		published,
	}
}
