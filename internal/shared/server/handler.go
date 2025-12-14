package server

import (
	"cli-t/internal/shared/file"
	"errors"
	"fmt"
	"io/fs"
	"path/filepath"
	"strings"
)

// mapPath converts request paths to file paths
// "/" → "index.html"
// "/foo.html" → "foo.html"
func mapPath(requestPath string) string {
	// TODO: You implement
	switch requestPath {
	case "/", "/index.html":
		return "index.html"
	default:
		return strings.TrimPrefix(requestPath, "/")
	}
}

/*
fullPath := filepath.Join(docRoot, requestPath)
// Result: "/etc/passwd"
// filepath.Join() already resolves  ..

cleaned := filepath.Clean(fullPath)
// Result: "/etc/passwd"
// Does NOT start with ".." - edge case
// But it's OUTSIDE /var/www!

another edge case noice
parentAbs := "/home/user/cli-t/www"
pathAbs := "/home/user/cli-t/www2/evil.html"
has prefix will fix this. will have to add slash?

Option 1

	completePath := filepath.Join(parent, relPath)
	completePath = filepath.Clean(completePath)

	parentAbs, err := filepath.Abs(parent)
	if err != nil {
		return "", err
	}

	pathAbs, err := filepath.Abs(completePath)
	if err != nil {
		return "", err
	}

	// Ensure parent ends with separator for proper prefix check
	if !strings.HasSuffix(parentAbs, string(filepath.Separator)) {
		parentAbs += string(filepath.Separator)
	}

	if !strings.HasPrefix(pathAbs, parentAbs) {
		//fmt.Errorf("path traversal detected")
		return "", fmt.Errorf("beta masti nahi")
	}

	return pathAbs, nil
*/
func safeJoin(parent, relPath string) (string, error) {
	// Get absolute paths
	parentAbs, err := filepath.Abs(parent)
	if err != nil {
		return "", err
	}

	// Join and clean
	joined := filepath.Join(parentAbs, relPath)
	cleanPath := filepath.Clean(joined)

	// Check if cleanPath is relative to parentAbs
	rel, err := filepath.Rel(parentAbs, cleanPath)
	if err != nil {
		return "", fmt.Errorf("path traversal detected: %w", err)
	}

	// If rel starts with "..", it's trying to escape!
	if strings.HasPrefix(rel, "..") {
		//fmt.Errorf("path traversal detected")
		return "", fmt.Errorf("beta masti nahi")
	}

	return cleanPath, nil
}

// serveFile reads and returns file content
// Returns: content, statusCode (200/404/500), error
func (s *Server) serveFile(requestPath string) ([]byte, int, error) {
	// TODO: You implement
	// Steps:
	// 1. Map path
	// 2. Join with docRoot (safe)
	// 3. Check if exists
	// 4. Read file
	// 5. Return content + status
	relPath := mapPath(requestPath)

	// Use filepath.Join() not string concat!
	//safe join is necessary since i could access
	//http://localhost:8000/../../lb/backend.go : Directory Traversal Attacks
	// completePath := filepath.Join(s.docRoot, relPath)
	completePath, err := safeJoin(s.docRoot, relPath)
	if err != nil {
		return []byte{}, 500, err // better status code
	}

	// later can think of streaming file over the network instead of one shot for big data
	data, err := file.ReadBytes(completePath)
	if err != nil {
		// different error for uncaught error as 500
		//if file does not exist 404
		// os.IsNotExist(err) old af
		if errors.Is(err, fs.ErrNotExist) {
			return []byte{}, 404, fmt.Errorf("file not found: %w", err)
		}
		return []byte{}, 500, fmt.Errorf("error reading file: %w", err)
	}

	return data, 200, nil
}
