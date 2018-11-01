package storage

import "time"

// Storage defines an interface for storing latest published time of the feed
type Storage interface {
	Update()
	GetLastPublishedTime(url string) time.Time
	PutLastPublishedTime(url string, lastPublished time.Time)
	Save()
}

type storage struct {
	config *Config
	state  map[string]time.Time
}

func (s storage) Update() {
	// TODO
}

func (s storage) GetLastPublishedTime(url string) time.Time {
	if v, ok := s.state[url]; ok {
		return v
	}
	return time.Unix(0, 0).UTC()
}

func (s storage) PutLastPublishedTime(url string, lastPublished time.Time) {
	s.state[url] = lastPublished
}

func (s storage) Save() {
	// TODO
}

// New constructs a new Storage
func New(config *Config) Storage {
	return storage{config: config}
}
