package lb

import (
	"cli-t/internal/shared/logger"
	"net/http"
	"net/url"
	"time"
)

type Checker struct {
	backends []*Backend
	interval time.Duration
	path     string
	timeout  time.Duration
	done     chan struct{}
}

func NewChecker(
	backends []*Backend,
	interval time.Duration,
	path string,
	timeout time.Duration,
	done chan struct{},
) *Checker {
	if done == nil {
		done = make(chan struct{})
	}

	return &Checker{
		backends: backends,
		interval: interval,
		path:     path,
		timeout:  timeout,
		done:     done,
	}
}

func (c *Checker) Start() {
	ticker := time.NewTicker(c.interval)
	go func() {
		for {
			select {
			case <-ticker.C: // ticker event is recived to channel
				c.checkAll()
			case <-c.done: // channel is closed
				ticker.Stop()
				return
			}
		}
	}()
}

// TODO: can check everything paralally
func (c *Checker) checkAll() {
	for _, backend := range c.backends {
		wasAlive := backend.IsAlive()
		alive := c.check(backend.URL)

		if wasAlive != alive {
			if alive {
				logger.Info("Backend recovered", "url", backend.URL)
			} else {
				logger.Warn("Backend unhealthy", "url", backend.URL)
			}
		} else if !alive {
			logger.Warn("Backend unalive", "url", backend.URL)
		}

		backend.SetAlive(alive)
	}
}

func (c *Checker) check(backendURL string) bool {
	// Create client with timeout
	client := &http.Client{
		Timeout: c.timeout,
	}

	base, err := url.Parse(backendURL)
	if err != nil {
		return false
	}

	base.Path = c.path

	resp, err := client.Get(base.String())
	if err != nil {
		return false
	}
	defer resp.Body.Close()

	return resp.StatusCode >= 200 && resp.StatusCode < 300
}

func (c *Checker) Stop() {
	close(c.done)
}
