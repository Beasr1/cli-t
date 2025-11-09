// internal/tools/lb/handler.go
package lb

import (
	"cli-t/internal/shared/logger"
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
}

// NewHandler creates a new load balancer handler
func NewHandler(backendURLs []string) (*Handler, error) {
	backends := make([]*Backend, 0, len(backendURLs))

	logger.Info("backend", "urls", backendURLs)

	for _, url := range backendURLs {
		backend, err := NewBackend(url)
		if err != nil {
			return nil, fmt.Errorf("invalid backend %s: %w", url, err)
		}
		backends = append(backends, backend)
	}

	return &Handler{
		backends:     backends,
		currentIndex: 0,
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

	current := h.currentIndex
	logger.Debug("current Index", "index", current)
	h.currentIndex = (h.currentIndex + 1) % len(h.backends) // can store length as well : will be lil faster
	return h.backends[current]
}

// Close cleans up resources
func (h *Handler) Close() error {
	// TODO: Close backend connections, stop health checks
	return nil
}
