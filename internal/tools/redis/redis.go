package redis

import (
	"cli-t/internal/command"
	"cli-t/internal/shared/logger"
	"cli-t/internal/tools/redis/server"
	"cli-t/internal/tools/redis/store"
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
	return "redis"
}

func (c *Command) Usage() string {
	return "redis [--host HOST] [--port PORT]"
}

func (c *Command) Description() string {
	return "Start a Redis-compatible in-memory key-value store server"
}

func (c *Command) ValidateArgs(args []string) error {
	return nil
}

// can also learn about detached mode?
// --db-file
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
			Default:   6379,
			Usage:     "port on which server is running",
		},
	}
}

/*
Redis Serialization Protocol (RESP)
*/
func (c *Command) Execute(ctx context.Context, args *command.Args) error {
	host, port := c.parseFlags(args.Flags)

	store := store.New()
	server := server.New(host, port, store)

	// Start server in goroutine
	go func() {
		// Start server (blocks until Ctrl+C) : This runs forever until interrupted
		if err := server.Start(); err != nil {
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
	return host, port
}
