package lb

import (
	"fmt"
	"net/http/httputil"
	"net/url"
	"sync"
)

// Backend represents a backend server
type Backend struct {
	URL   string
	Alive bool
	Proxy *httputil.ReverseProxy
	mu    *sync.RWMutex //  reads >> writes
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
		mu:    &sync.RWMutex{},
	}, nil
}

func (b *Backend) SetAlive(alive bool) {
	b.mu.Lock()
	defer b.mu.Unlock()
	b.Alive = alive
}

func (b *Backend) IsAlive() bool {
	b.mu.RLock()
	defer b.mu.RUnlock()
	return b.Alive
}
