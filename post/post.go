package post

import (
	"fmt"
	"html"
	"strings"
	"time"
)

// Post contains a title, a description and a url of a news
type Post struct {
	Title       string
	Description string
	URL         string
	Published   *time.Time
}

func clearString(s string) string {
	nospace := strings.TrimSpace(s)
	unescaped := html.UnescapeString(nospace)
	return unescaped
}

// New returns a pointer to the newly created Post struct
func New(title, description, url string, published *time.Time) *Post {
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

func (post Post) String() string {
	return fmt.Sprintf(
		"Post{\n\tTitle=%s\n\tDescription=%s\n\tURL=%s\n\tPublished=%v\n}",
		post.Title,
		post.Description,
		post.URL,
		post.Published,
	)
}
