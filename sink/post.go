package sink

import (
	"fmt"
	"html"
	"regexp"
	"strings"
	"time"
	"unicode/utf8"
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
		text = fmt.Sprintf("<b>%s</b>\n\n%s", p.title, p.url)
	default:
		text = fmt.Sprintf("<b>%s</b>\n\n%s\n\n%s", p.title, p.description, p.url)
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

func utf8Only(s string) string {
	result := ""
	for _, c := range s {
		if utf8.ValidRune(c) {
			result += string(c)
		}
	}
	return result
}

var htmlTagRegexp = regexp.MustCompile("<[^>]*>")

func clearString(s string) string {
	utf8Only := utf8Only(s)
	unescaped := html.UnescapeString(utf8Only)
	noCdata := strings.Replace(unescaped, "<![CDATA[", "", -1)
	noClosingCdata := strings.Replace(noCdata, "]]>", "", -1)
	noHTML := htmlTagRegexp.ReplaceAllString(noClosingCdata, "")
	trimmed := strings.TrimSpace(noHTML)
	return trimmed
}

// NewPost returns a pointer to the newly created Post struct
func NewPost(title, description, url string, published *time.Time) *Post {
	if published == nil {
		t := time.Unix(0, 0).UTC()
		published = &t
	}
	return &Post{
		limitLength(clearString(title), 100),
		limitLength(clearString(description), 280),
		clearString(url),
		*published,
	}
}

func limitLength(s string, limit int) string {
	if len(s) <= limit {
		return s
	}
	parts := strings.Split(s, " ")
	result := ""
	for i, word := range parts {
		if i != 0 {
			result += " "
		}
		result += word
		if len(result) >= limit {
			break
		}
	}
	return result + "…"
}
