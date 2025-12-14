package server

import (
	"cli-t/internal/shared/store/inmemory"
	"net"
)

type Server struct {
	host    string
	port    int
	store   inmemory.Store // Interface!
	baseURL string         // "http://localhost:8080"
}

func (s *Server) handleRequest(conn net.Conn) {
	// Parse HTTP request (you know this!)
	// Route based on method + path
	// Call appropriate handler

}
