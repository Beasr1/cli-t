package command

import (
	"context"
	"io"
)

// Command represents a CLI tool
type Command interface {
	// Metadata
	Name() string
	Usage() string
	Description() string

	// Execution
	Execute(ctx context.Context, args *Args) error

	// Validation
	ValidateArgs(args []string) error
}

// Args encapsulates all command inputs/outputs
type Args struct {
	// Command line arguments
	Positional []string
	Flags      map[string]interface{}

	// I/O
	Stdin  io.Reader
	Stdout io.Writer
	Stderr io.Writer

	// Environment
	Env    map[string]string
	Config *Config
}

// Config holds command-specific configuration
type Config struct {
	Verbose bool
	Debug   bool
	NoColor bool
	Output  string // output format
}
