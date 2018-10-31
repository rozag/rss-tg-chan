package config

import (
	"fmt"
)

// Config defines a single interface for processing command line flags for different parts of app
type Config interface {
	fmt.Stringer

	// RegisterFlags registers command line flags for the config
	RegisterFlags()
	// ValidateFlags validates command line flags for the config
	ValidateFlags() error
}
