package source

import (
	"fmt"
)

// Config contains params for the source of the rss feeds urls array
type Config struct {
	SourcesURL string
}

// LoadConfig returns a new shiny config
func LoadConfig(params map[string]string) *Config {
	return &Config{params["source"]}
}

// HelpLines returns a string slice with config format help lines
func (c *Config) HelpLines() []string {
	return []string{"source=SOURCE  // Sources json URL [Required]"}
}

// ValidateParams validates params for the config
func (c *Config) ValidateParams() error {
	if c.SourcesURL == "" {
		return fmt.Errorf("Sources json URL is required. Ensure providing it with the -source=URL flag")
	}
	return nil
}

func (c *Config) String() string {
	return fmt.Sprintf("{SourcesURL=%s}", c.SourcesURL)
}
