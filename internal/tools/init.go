package tools

import (
	"cli-t/internal/command"
	"cli-t/internal/shared/io"
	"cli-t/internal/tools/config"
	"cli-t/internal/tools/echo"
	"cli-t/internal/tools/list"
	versioncmd "cli-t/internal/tools/version"
	"fmt"
)

// Import all tools here
func init() {
	// Register all tools
	// Each tool is registered once when the package is imported

	// List of all tools to register
	tools := []command.Command{
		// built in
		config.New(),
		list.New(),
		versioncmd.New(),

		echo.New(),
	}

	// Register each tool
	for _, tool := range tools {
		if err := command.Register(tool); err != nil {
			// Panic during init is acceptable as it means
			// the program is misconfigured
			panic(fmt.Sprintf("failed to register %s command: %v",
				tool.Name(), err))
		}
	}

	// Log successful registration (optional)
	// fmt.Printf("Registered %d tools\n", len(tools))
}

// RegisterAll registers all available tools
func RegisterAll() error {
	// List of all tools
	tools := []command.Command{
		// built in
		config.New(),
		list.New(),
		versioncmd.New(),

		echo.New(),
		// wc.New(),
		// grep.New(),
		// cat.New(),
		// Add more tools here as you create them
	}

	// Register each tool
	for _, tool := range tools {
		if err := command.Register(tool); err != nil {
			io.Error("Failed to register tool",
				"tool", tool.Name(),
				"error", err)
			return err
		}
		io.Debug("Registered tool", "tool", tool.Name())
	}

	io.Info("All tools registered", "count", len(tools))
	return nil
}
