// internal/tools/lb/handler.go
package lb

import (
	"cli-t/internal/shared/logger"
	"fmt"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"
)

// Handler handles incoming HTTP requests and forwards them to backend servers
type Handler struct {
	// Future: backend pool, algorithm, health checker
	backend string
	proxy   *httputil.ReverseProxy
}

// NewHandler creates a new load balancer handler
func NewHandler(backendURL string) (*Handler, error) {
	// Parse URL
	targetURL, err := url.Parse(backendURL)
	if err != nil {
		return nil, fmt.Errorf("invalid backend URL: %w", err)
	}

	// Validate scheme (url.Parse doesn't error on missing scheme!)
	if targetURL.Scheme != "http" && targetURL.Scheme != "https" {
		return nil, fmt.Errorf("backend URL must start with http:// or https://")
	}

	// Also check if host is present
	if targetURL.Host == "" {
		return nil, fmt.Errorf("backend URL must include host")
	}

	// Create proxy
	proxy := httputil.NewSingleHostReverseProxy(targetURL)

	return &Handler{
		backend: backendURL,
		proxy:   proxy,
	}, nil
}

// ServeHTTP implements http.Handler interface
func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// Log request details
	logger.Info("Received request", "from", r.RemoteAddr)
	logger.Info("Request details",
		"method", r.Method,
		"path", r.URL.Path,
		"protocol", r.Proto,
	)

	// Log headers
	for key, values := range r.Header {
		logger.Info("Header", "name", key, "value", strings.Join(values, ", "))
	}

	// w.WriteHeader(http.StatusOK)
	// w.Write([]byte("Request received by load balancer\n"))
	h.proxy.ServeHTTP(w, r)
}

// Close cleans up resources
func (h *Handler) Close() error {
	// TODO: Close backend connections, stop health checks
	return nil
}
