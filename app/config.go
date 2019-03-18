package app

import (
	"fmt"
	"io/ioutil"
	"strconv"
	"strings"
	"time"

	"github.com/rozag/rss-tg-chan/sink"
	"github.com/rozag/rss-tg-chan/source"
	"github.com/rozag/rss-tg-chan/storage"
)

const (
	defaultLogLevel  = "d"
	defaultMinutes   = 60
	defaultWorkers   = 4
	defaultSingleRun = false
)

// Config is the main app config. It contains general params and configs for other parts of the app
type Config struct {
	LogLevel  string
	Period    time.Duration
	Workers   uint64
	SingleRun bool

	SourceConfig  *source.Config
	StorageConfig *storage.Config
	SinkConfig    *sink.Config
}

// LoadConfig returns a new shiny config
func LoadConfig(filename string) (*Config, error) {
	// Read config file content
	dat, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, fmt.Errorf("Failed to load the config file %s: %v", filename, err)
	}

	// Parse params from the config file into a map
	text := string(dat)
	lines := strings.Split(text, "\n")
	params := make(map[string]string)
	for _, line := range lines {
		clear := strings.TrimSpace(line)
		if clear == "" {
			continue
		}
		index := strings.Index(clear, "=")
		if index < 0 {
			continue
		}
		key := clear[:index]
		value := clear[index+1:]
		params[key] = value
	}

	// Get real params
	var logLevel string
	switch params["log"] {
	case "d":
		logLevel = "DEBUG"
	default:
		logLevel = "ERROR"
	}

	periodStr := params["period"]
	minutes, err := strconv.ParseInt(periodStr, 10, 64)
	var period time.Duration
	if err != nil {
		period = defaultMinutes * time.Minute
	} else {
		period = time.Duration(minutes) * time.Minute
		if period <= 0 {
			period = defaultMinutes * time.Minute
		}
	}

	workersStr := params["workers"]
	workers, err := strconv.ParseUint(workersStr, 10, 64)
	if err != nil {
		workers = defaultWorkers
	} else {
		if workers == 0 {
			workers = defaultWorkers
		}
	}

	singleRun := params["single"] == "true"

	return &Config{
		LogLevel:      logLevel,
		Period:        period,
		Workers:       workers,
		SingleRun:     singleRun,
		SourceConfig:  source.LoadConfig(params),
		StorageConfig: storage.LoadConfig(params),
		SinkConfig:    sink.LoadConfig(params),
	}, nil
}

// HelpLines returns a string slice with config format help lines
func (c *Config) HelpLines() []string {
	lines := []string{
		fmt.Sprintf("log=LOG_LEVEL  // Log level. \"e\" (for ERROR) or \"d\" (for DEBUG) log level. Default is %v", defaultLogLevel),
		fmt.Sprintf("period=PERIOD  // Period of the full load data and post results cycle in minutes. Default is %v", defaultPeriod),
		fmt.Sprintf("workers=WORKERS  // Number of workers for feeds processing. Default is %v", defaultWorkers),
		fmt.Sprintf("single=SINGLE  // If \"true\", only one load-and-post cycle will be executed. Default is %v", defaultSingleRun),
	}
	lines = append(lines, c.SourceConfig.HelpLines()...)
	lines = append(lines, c.StorageConfig.HelpLines()...)
	lines = append(lines, c.SinkConfig.HelpLines()...)
	return lines
}

// ValidateParams validates params for the config
func (c *Config) ValidateParams() error {
	if err := c.SourceConfig.ValidateParams(); err != nil {
		return err
	}
	if err := c.StorageConfig.ValidateParams(); err != nil {
		return err
	}
	if err := c.SinkConfig.ValidateParams(); err != nil {
		return err
	}
	return nil
}

// PrintDebugInfo logs values of it's params for debug purposes
func (c *Config) String() string {
	return fmt.Sprintf(
		"Config={LogLevel=%s Period=%v Workers=%d SingleRun=%v SourceConfig=%v StorageConfig=%v SinkConfig=%v}",
		c.LogLevel,
		c.Period,
		c.Workers,
		c.SingleRun,
		c.SourceConfig,
		c.StorageConfig,
		c.SinkConfig,
	)
}
