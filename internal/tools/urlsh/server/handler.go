package server

import "net"

// handler.go
func (s *Server) handlePOST(req *HTTPRequest, conn net.Conn) {
	// 1. Parse JSON body (research: json.Unmarshal)
	// 2. Validate URL
	// 3. Generate short code
	// 4. Save to storage
	// 5. Return JSON response
}

func (s *Server) handleGET(shortCode string, conn net.Conn) {
	// 1. Lookup in storage
	// 2. If found: 302 redirect
	// 3. If not: 404
}

func (s *Server) handleDELETE(shortCode string, conn net.Conn) {
	// 1. Delete from storage
	// 2. Return 200
}
