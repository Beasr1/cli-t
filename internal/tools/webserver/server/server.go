package server

import (
	"cli-t/internal/shared/logger"

	"context"
	"fmt"
	"net"
	"sync"
)

type Server struct {
	host string
	port int

	listener net.Listener      // TCP listener
	clients  map[net.Conn]bool // Active connections
	mu       sync.Mutex        // Protect clients map
	shutdown chan struct{}     // Signal to stop
}

func New(host string, port int) *Server {
	return &Server{
		host: host,
		port: port,

		clients:  make(map[net.Conn]bool),
		shutdown: make(chan struct{}),
	}
}

func (s *Server) Start(ctx context.Context) error {
	addr := fmt.Sprintf("%s:%d", s.host, s.port)

	// Listen on TCP (NOT http.Server!)
	listener, err := net.Listen("tcp", addr)
	if err != nil {
		return err
	}
	s.listener = listener

	logger.Info("Redis server listening", "addr", addr)

	// Accept connections loop
	for {
		_, err := listener.Accept()
		if err != nil {
			select {
			case <-s.shutdown:
				return nil // Clean shutdown
			default:
				logger.Error("Accept error", "error", err)
				continue
			}
		}

		// Handle each client in a goroutine
		// go s.handleClient(conn)
	}
}

func (s *Server) Stop(ctx context.Context) error {
	logger.Info("Shutting down server...")

	// 1. Stop accepting new connections
	close(s.shutdown)
	s.listener.Close()

	// 2. Wait for existing connections to finish (or timeout)
	done := make(chan struct{})

	go func() {
		// Close all active clients
		s.mu.Lock()
		for conn := range s.clients {
			conn.Close()
		}
		s.mu.Unlock()
		close(done)
	}()

	// 3. Wait for graceful shutdown or timeout
	select {
	case <-done:
		logger.Info("All clients closed gracefully")
		return nil
	case <-ctx.Done():
		logger.Warn("Shutdown timeout exceeded, forcing close")
		return ctx.Err()
	}
}
