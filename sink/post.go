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
	published   time.Time
}

// PublishedAfter returns true if the post was published after the specified time
func (p Post) PublishedAfter(t time.Time) bool {
	return p.published.After(t)
}

// YoungerThan returns true if this one post was published before the other one
func (p Post) YoungerThan(o Post) bool {
	return p.published.Before(o.published)
}

// GetPublished returns the published time of the post
func (p Post) GetPublished() time.Time {
	return p.published
}

// GetPublishableText returns ready for publishing post's text
func (p Post) GetPublishableText() string {
	if p.url == "" {
		return ""
	}
	var text string
	switch {
	case p.title == "" && p.description == "":
		text = p.url
	case p.title == "":
		text = fmt.Sprintf("%s\n\n%s", p.description, p.url)
	case p.description == "":
		text = fmt.Sprintf("*%s*\n\n%s", p.title, p.url)
	default:
		text = fmt.Sprintf("*%s*\n\n%s\n\n%s", p.title, p.description, p.url)
	}
	return text
}

func (p Post) String() string {
	return fmt.Sprintf(
		"Post={title=%s description=%s url=%s published=%v}",
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
		*published,
	}
}
