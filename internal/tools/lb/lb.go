// internal/tools/lb/lb.go
package lb

import (
	"cli-t/internal/command"
	"cli-t/internal/shared/logger"
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"
)

type Command struct{}

func New() command.Command {
	return &Command{}
}

func (c *Command) Name() string {
	return "lb"
}

func (c *Command) Usage() string {
	return "lb --port <port>"
}

func (c *Command) Description() string {
	return "Start an HTTP load balancer"
}

func (c *Command) ValidateArgs(args []string) error {
	return nil
}

func (c *Command) DefineFlags() []command.Flag {
	return []command.Flag{
		{
			Name:      "port",
			Shorthand: "p",
			Usage:     "Port to listen on",
			Type:      "int",
			Default:   8080,
		},
		{
			Name:      "backends",
			Shorthand: "b",
			Usage:     "Backend server URL (e.g., http://localhost:8081)",
			Type:      "string",
			Default:   "",
		},
		{
			Name:      "health-check-interval",
			Shorthand: "i",
			Usage:     "Health check interval (e.g., 10s, 1m)",
			Type:      "string",
			Default:   "10s",
		},
		{
			Name:      "health-check-path",
			Shorthand: "",
			Usage:     "Health check path",
			Type:      "string",
			Default:   "/",
		},
		{
			Name:      "health-check-timeout",
			Shorthand: "t",
			Usage:     "Health check timeout (e.g., 5s)",
			Type:      "string",
			Default:   "5s",
		},
	}
}

func (c *Command) Execute(ctx context.Context, args *command.Args) error {
	port, backends, healthCheckInterval, healthCheckPath, healthCheckTimeout := c.parseFlags(args.Flags)

	// Validate backend URL
	if len(backends) == 0 {
		return fmt.Errorf("backend URL is required")
	}

	// Create handler
	handler, err := NewHandler(backends, healthCheckInterval, healthCheckPath, healthCheckTimeout)
	if err != nil {
		return fmt.Errorf("failed to create handler: %w", err)
	}

	// Create server
	server := &http.Server{
		Addr:    fmt.Sprintf(":%d", port),
		Handler: handler,
	}

	// Start server in goroutine
	go func() {
		logger.Info("Starting load balancer", "port", port)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
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

	if err := handler.Close(); err != nil {
		logger.Error("Error closing handler", "error", err)
	}

	if err := server.Shutdown(shutdownCtx); err != nil {
		logger.Fatal("Server forced to shutdown", "error", err)
	}

	logger.Info("Server stopped gracefully")
	return nil
}

func (c *Command) parseFlags(flags map[string]interface{}) (int, []string, string, string, string) {
	port, _ := flags["port"].(int)
	backendsStr, _ := flags["backends"].(string)
	healthCheckInterval, _ := flags["health-check-interval"].(string)
	healthCheckPath, _ := flags["health-check-path"].(string)
	healthCheckTimeout, _ := flags["health-check-timeout"].(string)

	if backendsStr == "" {
		return port, []string{}, healthCheckInterval, healthCheckPath, healthCheckTimeout
	}

	parts := strings.Split(backendsStr, ",")
	backends := make([]string, 0, len(parts))

	//The range loop takes a snapshot of the slice at the start
	// hence if I keep appending parts : will not gie infinite loop
	for _, part := range parts {
		trimmed := strings.TrimSpace(part)
		if trimmed != "" {
			backends = append(backends, trimmed)
		}
	}

	logger.Debug("Flags processing",
		"port", port, "backends",
		backends, "healthCheckInterval",
		healthCheckInterval, "healthCheckPath",
		healthCheckPath, "healthCheckTimeout", healthCheckTimeout,
	)

	return port, backends, healthCheckInterval, healthCheckPath, healthCheckTimeout
}
