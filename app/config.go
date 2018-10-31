package app

import (
	"flag"
	"fmt"
	"time"

	"github.com/rozag/rss-tg-chan/sink"
	"github.com/rozag/rss-tg-chan/source"
	"github.com/rozag/rss-tg-chan/storage"
)

const (
	defaultLogLevel = "e"
	defaultPeriod   = time.Hour
	defaultWorkers  = 4
)

// Config is the main app config. It contains general params and configs for other parts of the app
type Config struct {
	LogLevel string
	Period   time.Duration
	Workers  uint

	logFlag     *string
	periodFlag  *time.Duration
	workersFlag *uint

	SourceConfig  *source.Config
	StorageConfig *storage.Config
	SinkConfig    *sink.Config
}

// NewConfig returns a new shiny config
func NewConfig() *Config {
	return &Config{
		SourceConfig:  source.NewConfig(),
		StorageConfig: storage.NewConfig(),
		SinkConfig:    sink.NewConfig(),
	}
}

// RegisterFlags registers command line flags for the config
func (c *Config) RegisterFlags() {
	c.logFlag = flag.String("log", defaultLogLevel, "Log level. \"e\" (for ERROR) or \"d\" (for DEBUG) log level")
	c.periodFlag = flag.Duration("period", defaultPeriod, "Period of the full load data and post results cycle")
	c.workersFlag = flag.Uint("workers", defaultWorkers, "Number of workers for feeds processing")

	c.SourceConfig.RegisterFlags()
	c.StorageConfig.RegisterFlags()
	c.SinkConfig.RegisterFlags()
}

// ValidateFlags validates command line flags for the config
func (c *Config) ValidateFlags() error {
	var logLevel string
	switch *c.logFlag {
	case "d":
		logLevel = "DEBUG"
	default:
		logLevel = "ERROR"
	}
	c.LogLevel = logLevel

	period := *c.periodFlag
	if period <= 0 {
		period = defaultPeriod
	}
	c.Period = period

	workers := *c.workersFlag
	if workers == 0 {
		workers = defaultWorkers
	}
	c.Workers = workers

	c.SourceConfig.ValidateFlags()
	c.StorageConfig.ValidateFlags()
	c.SinkConfig.ValidateFlags()

	return nil
}

// PrintDebugInfo logs values of it's params for debug purposes
func (c *Config) String() string {
	return fmt.Sprintf(
		"\nAppConfig:\n\tLogLevel='%s'\n\tPeriod=%v\n\tWorkers=%d\n\tSourceConfig=%v\n\tStorageConfig=%v\n\tSinkConfig=%v\n",
		c.LogLevel,
		c.Period,
		c.Workers,
		c.SourceConfig,
		c.StorageConfig,
		c.SinkConfig,
	)
}
