package main

import (
	"fmt"
	"log"
	"os"

	"github.com/hashicorp/logutils"
	"github.com/rozag/rss-tg-chan/app"
)

func printParseError(config *app.Config, err error) {
	fmt.Printf("[ERROR] Failed to parse config: %v\nConfig should have the following format:\n", err)
	lines := config.HelpLines()
	for _, line := range lines {
		fmt.Printf("\t%s\n", line)
	}
}

func main() {
	// Load config
	config, err := app.LoadConfig("config.ini")
	if err != nil {
		printParseError(config, err)
		return
	}
	err = config.ValidateParams()
	if err != nil {
		printParseError(config, err)
		return
	}

	// Set up the logs filter
	filter := logutils.LevelFilter{
		Levels:   []logutils.LogLevel{"DEBUG", "ERROR"},
		MinLevel: logutils.LogLevel(config.LogLevel),
		Writer:   os.Stderr,
	}
	log.SetOutput(&filter)

	// Print debug config info
	log.Printf("[DEBUG] %v", config)

	// Run the app
	app := app.New(config)
	done := make(chan bool)
	go app.Run(done)
	<-done
}
