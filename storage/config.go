package storage

import (
	"fmt"
)

// Config contains params for the storage
type Config struct {
	GithubToken        string
	GithubGistID       string
	GithubGistFileName string
}

// LoadConfig returns a new shiny config
func LoadConfig(params map[string]string) *Config {
	return &Config{params["githubToken"], params["githubGistID"], params["githubGistFileName"]}
}

// HelpLines returns a string slice with config format help lines
func (c *Config) HelpLines() []string {
	return []string{
		"githubToken=GITHUB_TOKEN  // Github personal access token with the 'gist' scope [Required]",
		"githubGistID=GITHUB_GIST_ID  // Id of a gist used as a storage [Required]",
		"githubGistFileName=GITHUB_GIST_FILE_NAME  // Name of a file in the gist, which is an actual storage [Required]",
	}
}

// ValidateParams validates params for the config
func (c *Config) ValidateParams() error {
	if c.GithubToken == "" {
		return fmt.Errorf("Github token is required. Ensure providing it with the -githubToken=TOKEN flag")
	}
	if c.GithubGistID == "" {
		return fmt.Errorf("Github gist id is required. Ensure providing it with the -githubGistID=ID flag")
	}
	if c.GithubGistFileName == "" {
		return fmt.Errorf("Github gist file name is required. Ensure providing it with the -githubGistFileName=FILENAME flag")
	}
	return nil
}

func (c *Config) String() string {
	return fmt.Sprintf("{GithubGistID=%s GithubGistFileName=%s}", c.GithubGistID, c.GithubGistFileName)
}
