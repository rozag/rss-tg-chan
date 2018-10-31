package sink

import (
	"flag"
	"fmt"
)

// Config contains params for the sink of the news from rss feeds
type Config struct {
	TgBotToken string
	TgChannel  string

	tgTokenFlag   *string
	tgChannelFlag *string
}

// NewConfig returns a new shiny config
func NewConfig() *Config {
	return &Config{}
}

// RegisterFlags registers command line flags for the config
func (c *Config) RegisterFlags() {
	c.tgTokenFlag = flag.String("tgToken", "", "Telegram bot token (get it here https://t.me/BotFather) [Required]")
	c.tgChannelFlag = flag.String("tgChannel", "", "Telegram channel id (e.g. @mynewchannel) [Required]")
}

// ValidateFlags validates command line flags for the config
func (c *Config) ValidateFlags() error {
	tgBotToken := *c.tgTokenFlag
	if tgBotToken == "" {
		return fmt.Errorf("Telegram bot token is required. Ensure providing it with the -tgToken=BOT_TOKEN flag")
	}
	c.TgBotToken = tgBotToken

	tgChannel := *c.tgChannelFlag
	if tgChannel == "" {
		return fmt.Errorf("Telegram channel id is required. Ensure providing it with the -tgChannel=CHANNEL_ID flag")
	}
	c.TgChannel = tgChannel

	return nil
}

func (c *Config) String() string {
	return fmt.Sprintf("(TgChannel='%s')", c.TgChannel)
}
