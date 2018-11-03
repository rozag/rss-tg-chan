package storage

import (
	"flag"
	"fmt"
)

// Config contains params for the storage
type Config struct {
	GithubToken        string
	GithubGistID       string
	GithubGistFileName string

	githubTokenFlag        *string
	githubGistIDFlag       *string
	githubGistFileNameFlag *string
}

// NewConfig returns a new shiny config
func NewConfig() *Config {
	return &Config{}
}

// RegisterFlags registers command line flags for the config
func (c *Config) RegisterFlags() {
	c.githubTokenFlag = flag.String("githubToken", "", "Github personal access token with the 'gist' scope [Required]")
	c.githubGistIDFlag = flag.String("githubGistID", "", "Id of a gist used as a storage [Required]")
	c.githubGistFileNameFlag = flag.String("githubGistFileName", "", "Name of a file in the gist, which is an actual storage [Required]")
}

// ValidateFlags validates command line flags for the config
func (c *Config) ValidateFlags() error {
	githubToken := *c.githubTokenFlag
	if githubToken == "" {
		return fmt.Errorf("Github token is required. Ensure providing it with the -githubToken=TOKEN flag")
	}
	c.GithubToken = githubToken

	githubGistID := *c.githubGistIDFlag
	if githubGistID == "" {
		return fmt.Errorf("Github gist id is required. Ensure providing it with the -githubGistID=ID flag")
	}
	c.GithubGistID = githubGistID

	githubGistFileName := *c.githubGistFileNameFlag
	if githubGistFileName == "" {
		return fmt.Errorf("Github gist file name is required. Ensure providing it with the -githubGistFileName=FILENAME flag")
	}
	c.GithubGistFileName = githubGistFileName

	return nil
}

func (c *Config) String() string {
	return fmt.Sprintf("{GithubGistID=%s GithubGistFileName=%s}", c.GithubGistID, c.GithubGistFileName)
}
