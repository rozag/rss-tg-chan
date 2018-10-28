package post

import (
	"fmt"
)

// Post contains a title, a description and a url of a news
type Post struct {
	Title       string
	Description string
	URL         string
}

// New returns a pointer to the newly created Post struct
func New(title, description, url string) *Post {
	return &Post{title, description, url}
}

func (post Post) String() string {
	return fmt.Sprintf("Post{\n\tTitle=%s\n\tDescription=%s\n\tURL=%s\n\n}", post.Title, post.Description, post.URL)
}
