package main

import (
	"flag"
	"log"
	"os"

	"github.com/hashicorp/logutils"
	"github.com/rozag/rss-tg-chan/app"
)

func main() {
	// Register, parse and validate flags
	config := app.NewConfig()
	config.RegisterFlags()
	flag.Parse()
	err := config.ValidateFlags()
	if err != nil {

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
