package server

import (
	"bufio"
	"fmt"
	"strings"
)

type Header struct {
	Key   string
	Value string
}

type HTTPRequest struct {
	// What do you need from the HTTP request?
	// Method? Path? Version?
	Method  string
	Path    string
	Version string
	Headers []Header
}

func ParseRequest(reader *bufio.Reader) (*HTTPRequest, error) {
	// 1. Read first line → parse method/path/version
	line, err := reader.ReadString('\n')
	if err != nil {
		return nil, fmt.Errorf("failed to read request line: %w", err)
	}

	line = strings.TrimRight(line, "\r\n")
	components := strings.Split(line, " ")

	if len(components) != 3 {
		return nil, fmt.Errorf("invalid request line format: expected 3 parts, got %d", len(components))
	}

	// 2. Loop: read headers until empty line
	//    (don't parse them, just consume)
	headers := []Header{}
	for {
		header, err := reader.ReadString('\n')
		if err != nil {
			return nil, fmt.Errorf("failed to read request line: %w", err)
		}

		header = strings.TrimRight(header, "\r\n")
		if header == "" { // Empty after trim
			break
		}

		headerComp := strings.SplitN(header, ":", 2) // Date: Mon, 27 Jul 2009 12:28:53 GMT
		if len(headerComp) != 2 {
			return nil, fmt.Errorf("invalid request line format: expected 2 parts, got %d", len(headerComp))
		}

		headers = append(headers, Header{
			Key:   headerComp[0],
			Value: strings.TrimSpace(headerComp[1]),
		})
	}

	httpReq := HTTPRequest{
		Method:  components[0],
		Path:    components[1],
		Version: components[2],
		Headers: headers,
	}

	// 3. Return HTTPRequest with method/path/version
	return &httpReq, nil
}
