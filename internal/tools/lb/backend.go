package lb

import (
	"fmt"
	"net/http/httputil"
	"net/url"
)

// Backend represents a backend server
type Backend struct {
	URL   string
	Alive bool
	Proxy *httputil.ReverseProxy
}

// NewBackend creates a new backend server
func NewBackend(backendURL string) (*Backend, error) {
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

	return &Backend{
		URL:   backendURL,
		Proxy: proxy,
		Alive: true, //assume all backends are healthy until proven otherwise.
	}, nil
}
