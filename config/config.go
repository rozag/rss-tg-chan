package config

import (
	"flag"
	"fmt"
	"time"
)

// Config contains global params
type Config struct {
	LogLevel   string
	SourcesURL string
	Period     time.Duration
	Workers    uint
	TgBotToken string
	TgChannel  string
}

// ParseFlags parses the command line arguments and returns Config or error
func ParseFlags() (*Config, error) {
	logLevelFlag := flag.String("log", "e", "Log level. \"e\" (for ERROR) or \"d\" (for DEBUG) log level")
	sourcesURLFlag := flag.String("source", "", "Sources json URL. Required")
	periodFlag := flag.Duration("period", time.Hour, "Period of the full load data and post results cycle")
	workersFlag := flag.Uint("workers", 4, "Number of workers for feeds processing")
	tokenFlag := flag.String("tgToken", "", "Telegram bot token (get it here https://t.me/BotFather). Required")
	channelFlag := flag.String("tgChannel", "", "Telegram channel id (e.g. @mynewchannel). Required")
	flag.Parse()

	var logLevel string
	switch *logLevelFlag {
	case "d":
		logLevel = "DEBUG"
	default:
		logLevel = "ERROR"
	}

	sourcesURL := *sourcesURLFlag
	if sourcesURL == "" {
		return nil, fmt.Errorf("Sources json URL is required. Ensure providing it with the -source=URL flag")
	}

	tgBotToken := *tokenFlag
	if tgBotToken == "" {
		return nil, fmt.Errorf("Telegram bot token is required. Ensure providing it with the -tgToken=BOT_TOKEN flag")
	}

	tgChannel := *channelFlag
	if tgChannel == "" {
		return nil, fmt.Errorf("Telegram channel id is required. Ensure providing it with the -tgChannel=CHANNEL_ID flag")
	}

	return &Config{
		LogLevel:   logLevel,
		SourcesURL: sourcesURL,
		Period:     *periodFlag,
		Workers:    *workersFlag,
		TgBotToken: tgBotToken,
		TgChannel:  tgChannel,
	}, nil
}
