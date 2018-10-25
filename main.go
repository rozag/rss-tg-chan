package main

import (
	"flag"
	"log"
	"os"

	"github.com/hashicorp/logutils"
)

func main() {
	logLevelFlag := flag.String("log", "e", "Log level. \"e\" (for ERROR) or \"d\" (for DEBUG) log level")
	sourcesURLFlag := flag.String("source", "", "Sources json URL. Required.")
	flag.Parse()

	var logLevel string
	switch *logLevelFlag {
	case "d":
		logLevel = "DEBUG"
	default:
		logLevel = "ERROR"
	}
	filter := &logutils.LevelFilter{
		Levels:   []logutils.LogLevel{"DEBUG", "ERROR"},
		MinLevel: logutils.LogLevel(logLevel),
		Writer:   os.Stderr,
	}
	log.SetOutput(filter)
	log.Printf("[DEBUG] Using %s log level", logLevel)

	sourcesURL := *sourcesURLFlag
	if sourcesURL == "" {
		log.Println("[ERROR] Sources json URL is required. Ensure providing it with the -source=URL flag")
		return
	}
	log.Printf("[DEBUG] Using %s sources url", sourcesURL)

	// TODO
}
