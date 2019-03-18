package config

import (
	"fmt"
)

// Config defines a single interface for processing params for different parts of app
type Config interface {
	fmt.Stringer

	// HelpLines returns a string slice with config format help lines
	HelpLines() []string
	// ValidateParams validates params for the config
	ValidateParams() error
}
