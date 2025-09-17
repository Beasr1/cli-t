package main

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"time"

	"cli-t/internal/command"
	"cli-t/internal/config"
	"cli-t/internal/shared/io"
	_ "cli-t/internal/tools" // Register tools
	"cli-t/pkg/version"

	"github.com/spf13/cobra"
)

var (
	// Global flags
	verbose  bool
	debug    bool
	trace    bool
	noColor  bool
	output   string
	logLevel string
	logFile  string
)

func main() {
	rootCmd := &cobra.Command{
		Use:   "cli-t",
		Short: "The Developer's Swiss Army Knife",
		Long: `CLI-T consolidates essential developer tools into one powerful, customizable CLI.
No more installing dozens of tools - get everything from wc to redis-server in a single binary.`,
		Version: version.String(),
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			// Determine log level from flags
			level := determineLogLevel()

			// Initialize logger
			logConfig := &io.LogConfig{
				Level:      level,
				Format:     "console", // Could add --log-format flag later
				Output:     "stderr",
				NoColor:    noColor,
				ShowCaller: debug || trace,
				ShowTime:   true,
			}

			// Override output if log file specified
			if logFile != "" {
				logConfig.Output = logFile
			}

			if err := io.Initialize(logConfig); err != nil {
				return fmt.Errorf("failed to initialize logger: %w", err)
			}

			// Log startup info
			io.Debug("CLI-T starting",
				"version", version.Version,
				"command", cmd.Name(),
				"args", args,
				"pid", os.Getpid(),
			)

			// Load configuration
			if _, err := config.Load(); err != nil {
				io.Warn("Failed to load config", "error", err)
				// Don't fail, just use defaults
			}

			return nil
		},
		PersistentPostRun: func(cmd *cobra.Command, args []string) {
			// Ensure logs are flushed
			io.Sync()
		},
	}

	// Global flags
	rootCmd.PersistentFlags().BoolVarP(&verbose, "verbose", "v", false, "Verbose output")
	rootCmd.PersistentFlags().BoolVar(&debug, "debug", false, "Debug mode")
	rootCmd.PersistentFlags().BoolVar(&trace, "trace", false, "Trace mode (very verbose)")
	rootCmd.PersistentFlags().BoolVar(&noColor, "no-color", false, "Disable colored output")
	rootCmd.PersistentFlags().StringVarP(&output, "output", "o", "plain", "Output format (plain, json, yaml)")
	rootCmd.PersistentFlags().StringVar(&logLevel, "log-level", "", "Log level (trace, verbose, debug, info, warn, error)")
	rootCmd.PersistentFlags().StringVar(&logFile, "log-file", "", "Write logs to file instead of stderr")

	// Mark some flags as hidden (advanced users only)
	rootCmd.PersistentFlags().MarkHidden("trace")
	rootCmd.PersistentFlags().MarkHidden("log-level")

	// Add subcommands from registry
	for name, cmd := range command.GetAll() {
		rootCmd.AddCommand(createSubcommand(name, cmd))
	}

	// Add built-in commands
	rootCmd.AddCommand(createConfigCommand())
	rootCmd.AddCommand(createListCommand())
	rootCmd.AddCommand(createVersionCommand())

	// Execute
	if err := rootCmd.Execute(); err != nil {
		io.Error("Command failed", "error", err)
		os.Exit(1)
	}
}

// determineLogLevel figures out the log level from various flags
func determineLogLevel() string {
	// Explicit log level takes precedence
	if logLevel != "" {
		return logLevel
	}

	// Then check boolean flags
	if trace {
		return "trace"
	}
	if debug {
		return "debug"
	}
	if verbose {
		return "verbose"
	}

	// Check environment variables
	if level := os.Getenv("CLI_T_LOG_LEVEL"); level != "" {
		return level
	}
	if os.Getenv("CLI_T_DEBUG") == "true" {
		return "debug"
	}

	// Default
	return "info"
}

func createSubcommand(name string, cmd command.Command) *cobra.Command {
	return &cobra.Command{
		Use:   name + " " + cmd.Usage(),
		Short: cmd.Description(),
		Long:  cmd.Description(), // Could be enhanced with long descriptions
		Args:  cobra.ArbitraryArgs,
		RunE: func(cobraCmd *cobra.Command, args []string) error {
			// Log command execution
			start := time.Now()
			io.LogCommand(name, args)

			// Create command args
			cmdArgs := &command.Args{
				Positional: args,
				Flags:      make(map[string]interface{}),
				Stdin:      os.Stdin,
				Stdout:     os.Stdout,
				Stderr:     os.Stderr,
				Env:        envToMap(os.Environ()),
				Config: &command.Config{
					Verbose: verbose || debug || trace,
					Debug:   debug || trace,
					NoColor: noColor,
					Output:  output,
				},
			}

			// TODO: Parse command-specific flags from cobraCmd

			// Execute command
			err := cmd.Execute(cobraCmd.Context(), cmdArgs)

			// Log completion
			io.LogDuration(name, start, "success", err == nil)

			return err
		},
	}
}

func createConfigCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "config",
		Short: "Manage CLI-T configuration",
		Long:  "View and edit CLI-T configuration settings",
	}

	// Subcommands
	cmd.AddCommand(&cobra.Command{
		Use:   "show",
		Short: "Show current configuration",
		RunE: func(cmd *cobra.Command, args []string) error {
			cfg, err := config.Load()
			if err != nil {
				return err
			}

			// TODO: Pretty print config
			fmt.Printf("%+v\n", cfg)
			return nil
		},
	})

	cmd.AddCommand(&cobra.Command{
		Use:   "edit",
		Short: "Edit configuration file",
		RunE: func(cmd *cobra.Command, args []string) error {
			// TODO: Open config in editor
			return fmt.Errorf("not implemented yet")
		},
	})

	return cmd
}

func createListCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "list",
		Short: "List available tools",
		Long:  "Display all available CLI-T tools",
		Run: func(cmd *cobra.Command, args []string) {
			tools := command.List()

			io.Info("Available tools", "count", len(tools))

			fmt.Println("\nAvailable tools:")
			fmt.Println(strings.Repeat("-", 50))

			for _, name := range tools {
				if cmd, ok := command.Get(name); ok {
					fmt.Printf("  %-15s %s\n", name, cmd.Description())
				}
			}

			fmt.Printf("\nTotal: %d tools\n", len(tools))
			fmt.Println("\nUse 'cli-t <tool> --help' for tool-specific help")
		},
	}
}

func createVersionCommand() *cobra.Command {
	var (
		short   bool
		jsonOut bool
	)

	cmd := &cobra.Command{
		Use:   "version",
		Short: "Show version information",
		Long:  "Display detailed version information about CLI-T",
		Run: func(cmd *cobra.Command, args []string) {
			if short {
				fmt.Println(version.Version)
				return
			}

			if jsonOut {
				data, _ := json.MarshalIndent(version.GetInfo(), "", "  ")
				fmt.Println(string(data))
				return
			}

			// Default detailed output
			fmt.Println(version.DetailedString())
		},
	}

	cmd.Flags().BoolVarP(&short, "short", "s", false, "Show version number only")
	cmd.Flags().BoolVar(&jsonOut, "json", false, "Output version info as JSON")

	return cmd
}

func envToMap(env []string) map[string]string {
	m := make(map[string]string)
	for _, e := range env {
		if i := strings.Index(e, "="); i >= 0 {
			m[e[:i]] = e[i+1:]
		}
	}
	return m
}
