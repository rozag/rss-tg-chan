package sink

import (
	"fmt"
)

// Config contains params for the sink of the news from rss feeds
type Config struct {
	TgBotToken string
	TgChannel  string
}

// LoadConfig returns a new shiny config
func LoadConfig(params map[string]string) *Config {
	return &Config{params["tgToken"], params["tgChannel"]}
}

// HelpLines returns a string slice with config format help lines
func (c *Config) HelpLines() []string {
	return []string{
		"tgToken=TG_TOKEN  // Telegram bot token (get it here https://t.me/BotFather) [Required]",
		"tgChannel=TG_CHANNEL  // Telegram channel id (e.g. @mynewchannel) [Required]",
	}
}

// ValidateParams validates params for the config
func (c *Config) ValidateParams() error {
	if c.TgBotToken == "" {
		return fmt.Errorf("Telegram bot token is required. Ensure providing it with the -tgToken=BOT_TOKEN flag")
	}
	if c.TgChannel == "" {
		return fmt.Errorf("Telegram channel id is required. Ensure providing it with the -tgChannel=CHANNEL_ID flag")
	}
	return nil
}

func (c *Config) String() string {
	return fmt.Sprintf("{TgChannel=%s}", c.TgChannel)
}
