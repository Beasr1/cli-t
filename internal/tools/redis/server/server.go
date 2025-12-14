package server

import (
	"cli-t/internal/shared/logger"
	inmemory "cli-t/internal/shared/store/inmemory"
	"cli-t/internal/tools/redis/protocol"

	"context"
	"fmt"
	"io"
	"net"
	"sync"
)

type Server struct {
	host  string
	port  int
	store inmemory.Store

	listener net.Listener      // TCP listener
	clients  map[net.Conn]bool // Active connections
	mu       sync.Mutex        // Protect clients map
	shutdown chan struct{}     // Signal to stop
}

func New(host string, port int, store inmemory.Store) *Server {
	return &Server{
		host:     host,
		port:     port,
		store:    store,
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

	// Start background expiry worker
	go s.store.StartExpiryWorker(ctx)
	logger.Info("Started expiry cleanup worker")

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
		go s.handleClient(conn)
	}
}

func (s *Server) Stop(ctx context.Context) error { // ← Return error
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

// Make sure connection is persistent so tcp handshake is not happening on every hit
// TODO: can set a ideal timeout
func (s *Server) handleClient(conn net.Conn) {
	// Track this client
	s.mu.Lock()
	s.clients[conn] = true
	s.mu.Unlock()

	remoteAddr := conn.RemoteAddr().String()
	logger.Info("Client connected", "addr", remoteAddr)

	defer func() {
		s.closeClient(conn)
		logger.Info("Client disconnected", "addr", remoteAddr)
	}()

	// Buffer gets overwritten each time : fast reuse
	/*
		Known issue for now :
		Partial Reads (Message Split Across Packets) : one command in multiple packet
		Multiple Commands in One Read : multiple command in one packet
		we will test them then fix this issue
	*/
	buffer := make([]byte, 4096)

	for { // ← Keep reading until client disconnects!
		n, err := conn.Read(buffer)
		if err != nil {
			if err == io.EOF {
				logger.Debug("Client disconnected")
			} else {
				logger.Error("Read error", "error", err)
			}
			return // Exit loop, defer calls Close()
		}

		// Parse command
		data := buffer[:n]
		msg, _, err := protocol.Parse(data)
		if err != nil {
			conn.Write([]byte("-ERR protocol error\r\n"))
			continue // ← Stay in loop, don't close!
		}

		// Handle command
		response := s.handleCommand(msg)
		conn.Write(response.Serialize())

		// ← Loop back, read next command from SAME client
	}
}

func (s *Server) closeClient(conn net.Conn) {
	conn.Close()
	s.mu.Lock()
	delete(s.clients, conn)
	s.mu.Unlock()
}
