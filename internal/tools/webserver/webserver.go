// internal/tools/webserver/webserver.go
package webserver

import (
	"cli-t/internal/command"
	"cli-t/internal/shared/logger"
	"cli-t/internal/tools/webserver/server"

	"context"
	"os"
	"os/signal"
	"syscall"
	"time"
)

type Command struct{}

func New() command.Command {
	return &Command{}
}

func (c *Command) Name() string {
	return "server"
}

func (c *Command) Usage() string {
	return "server --port <port>"
}

func (c *Command) Description() string {
	return "Start an HTTP server"
}

func (c *Command) ValidateArgs(args []string) error {
	return nil
}

func (c *Command) DefineFlags() []command.Flag {
	return []command.Flag{
		{
			Name:      "host",
			Shorthand: "",
			Type:      "string",
			Default:   "127.0.0.1",
			Usage:     "host on which server is running",
		},
		{
			Name:      "port",
			Shorthand: "p",
			Type:      "int",
			Default:   8000,
			Usage:     "port on which server is running",
		},
	}
}

func (c *Command) Execute(ctx context.Context, args *command.Args) error {
	host, port := c.parseFlags(args.Flags)

	// Create server
	server := server.New(host, port)

	// Start server in goroutine
	go func() {
		// Start server (blocks until Ctrl+C) : This runs forever until interrupted
		if err := server.Start(ctx); err != nil {
			// make sure ctrl C : does not return server failed
			logger.Fatal("Server failed", "error", err)
		}
	}()

	// Wait for interrupt signal
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	logger.Info("Shutting down server...")

	// Graceful shutdown
	shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := server.Stop(shutdownCtx); err != nil {
		logger.Fatal("Server forced to shutdown", "error", err)
	}

	logger.Info("Server stopped gracefully")
	return nil
}

func (c *Command) parseFlags(flags map[string]interface{}) (string, int) {
	host, _ := flags["host"].(string)
	port, _ := flags["port"].(int)

	logger.Debug("Flags processing",
		"port", port, "backends",
	)

	return host, port
}
