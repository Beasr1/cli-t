// internal/tools/lb/handler.go
package lb

import (
	"cli-t/internal/shared/logger"
	"time"

	"fmt"
	"net/http"
	"sync"
)

// Handler handles incoming HTTP requests and forwards them to backend servers
type Handler struct {
	// Future: backend pool, algorithm, health checker
	backends     []*Backend
	currentIndex int
	mu           sync.Mutex
	healthCheck  *Checker
}

// NewHandler creates a new load balancer handler
func NewHandler(backendURLs []string, healthCheckInterval, healthCheckPath, healthCheckTimeout string) (*Handler, error) {
	backends := make([]*Backend, 0, len(backendURLs))

	logger.Info("backend", "urls", backendURLs)

	for _, url := range backendURLs {
		backend, err := NewBackend(url)
		if err != nil {
			return nil, fmt.Errorf("invalid backend %s: %w", url, err)
		}
		backends = append(backends, backend)
	}

	// Parse duration strings
	interval, err := time.ParseDuration(healthCheckInterval)
	if err != nil {
		return nil, fmt.Errorf("invalid health check interval: %w", err)
	}

	timeout, err := time.ParseDuration(healthCheckTimeout)
	if err != nil {
		return nil, fmt.Errorf("invalid health check timeout: %w", err)
	}

	checker := NewChecker(backends, interval, healthCheckPath, timeout, nil)
	checker.Start()

	return &Handler{
		backends:     backends,
		currentIndex: 0,
		healthCheck:  checker,
	}, nil
}

// ServeHTTP implements http.Handler interface
func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// Get next backend using round robin
	backend := h.nextBackend()

	// Log request details
	logger.Info("Forwarding request",
		"from", r.RemoteAddr,
		"to", backend.URL,
		"path", r.URL.Path,
		"method", r.Method,
		"protocol", r.Proto,
	)

	backend.Proxy.ServeHTTP(w, r)
}

// Simple Round Robin
func (h *Handler) nextBackend() *Backend {
	h.mu.Lock()
	defer h.mu.Unlock()

	for i := 0; i < len(h.backends); i++ {
		backend := h.backends[h.currentIndex]
		h.currentIndex = (h.currentIndex + 1) % len(h.backends) // can store length as well : will be lil faster

		logger.Info("index", "index", h.currentIndex)
		if backend.IsAlive() {
			return backend
		}
	}

	// All dead, return first anyway
	return h.backends[0]
}

// Close cleans up resources
func (h *Handler) Close() error {
	// TODO: Close backend connections, stop health checks
	h.healthCheck.Stop()
	return nil
}
