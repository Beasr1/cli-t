package server

import (
	"bufio"
	"cli-t/internal/shared/logger"

	"context"
	"fmt"
	"net"
	"sync"
)

type Server struct {
	host string
	port int

	docRoot string

	listener net.Listener      // TCP listener
	clients  map[net.Conn]bool // Active connections
	mu       sync.Mutex        // Protect clients map
	shutdown chan struct{}     // Signal to stop
}

func New(host string, port int, docRoot string) *Server {
	return &Server{
		host: host,
		port: port,

		clients:  make(map[net.Conn]bool),
		shutdown: make(chan struct{}),

		docRoot: docRoot,
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

	logger.Info("Server listening", "addr", addr)

	// Accept connections loop
	for {
		conn, err := listener.Accept()
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
		go s.handleConnection(conn)
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
			delete(s.clients, conn)
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

func (s *Server) handleConnection(conn net.Conn) {
	defer func() {
		// Cleanup code here
		// This runs when function exits, no matter how
		s.mu.Lock()
		delete(s.clients, conn)
		s.mu.Unlock()
		conn.Close()
	}()

	s.mu.Lock()
	s.clients[conn] = true
	s.mu.Unlock()

	// TEMPORARY: See concurrency in action
	// logger.Info("Connection started", "addr", conn.RemoteAddr())
	// time.Sleep(5 * time.Second) // 5 second delay
	// logger.Info("Connection processing", "addr", conn.RemoteAddr())

	reader := bufio.NewReader(conn)
	req, err := ParseRequest(reader)
	if err != nil {
		// Send what string to conn
		// HTTP/1.1 400 Bad Request\r\n\r\nInvalid request\r\n
		response := "HTTP/1.1 400 Bad Request\r\n\r\nInvalid request\r\n"
		conn.Write([]byte(response)) // Convert string to []byte
		return
	}

	// Serve file
	content, statusCode, err := s.serveFile(req.Path)
	if err != nil {
		// handled below
		logger.Error("error", "err", err)
	}

	switch statusCode {
	case 200:
		// Send header then raw file content
		header := "HTTP/1.1 200 OK\r\n\r\n"
		conn.Write([]byte(header))
		conn.Write(content)

	case 404:
		response := "HTTP/1.1 404 Not Found\r\n\r\nFile not found\r\n"
		conn.Write([]byte(response))

	case 500:
		response := "HTTP/1.1 500 Internal Server Error\r\n\r\nServer error\r\n"
		conn.Write([]byte(response))
	}

}
