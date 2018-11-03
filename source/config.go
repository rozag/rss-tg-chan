package source

import (
	"flag"
	"fmt"
)

// Config contains params for the source of the rss feeds urls array
type Config struct {
	SourcesURL  string
	sourcesFlag *string
}

// NewConfig returns a new shiny config
func NewConfig() *Config {
	return &Config{}
}

// RegisterFlags registers command line flags for the config
func (c *Config) RegisterFlags() {
	c.sourcesFlag = flag.String("source", "", "Sources json URL [Required]")
}

// ValidateFlags validates command line flags for the config
func (c *Config) ValidateFlags() error {
	sourcesURL := *c.sourcesFlag
	if sourcesURL == "" {
		return fmt.Errorf("Sources json URL is required. Ensure providing it with the -source=URL flag")
	}

	c.SourcesURL = sourcesURL
	return nil
}

func (c *Config) String() string {
	return fmt.Sprintf("{SourcesURL=%s}", c.SourcesURL)
}
