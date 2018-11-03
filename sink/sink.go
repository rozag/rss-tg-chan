package sink

// Sink defines an interface for publishing posts from feeds
type Sink struct {
	config *Config
}

// New constructs a new Sink
func New(config *Config) *Sink {
	return &Sink{config}
}
