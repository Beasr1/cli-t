package server

import (
	"bufio"
	"fmt"
	"io"
	"strconv"
	"strings"
)

// Can have better mapping of keys for O(1) search if required
type Header struct {
	Key   string
	Value string
}

type HTTPRequest struct {
	// What do you need from the HTTP request?
	// Method? Path? Version?
	Method      string
	Path        string
	Version     string
	Headers     []Header
	RequestBody []byte
}

func findHeaderByKey(headers []Header, key string) (Header, bool) {
	for _, header := range headers { // Iterate over the slice
		if header.Key == key { // Check the specific field
			return header, true // Return the struct and true if found
		}
	}
	return Header{}, false // Return an empty struct and false if not found
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

	/*
		POST / HTTP/1.1\r\n
		Content-Type: application/json\r\n
		Content-Length: 34\r\n
		\r\n
		{"url":"https://example.com"}  // ← Body here!
	*/

	buff, ok := findHeaderByKey(headers, "Content-Length") // should it by case insensitive
	if !ok {
		// if get ignore ?
		return nil, fmt.Errorf("invalid request line format: Content-Length. recived none")
	}

	bodyLength, err := strconv.ParseInt(buff.Value, 10, 64)
	if err != nil {
		// if get ignore ?
		return nil, fmt.Errorf("invalid request line format: Content-Length. recived none")
	}

	// Based on content-type it can be xml, json , plain text, etc
	requestBody := make([]byte, bodyLength)
	io.ReadFull(reader, requestBody)

	httpReq := HTTPRequest{
		Method:      components[0],
		Path:        components[1],
		Version:     components[2],
		Headers:     headers,
		RequestBody: requestBody,
	}

	// 3. Return HTTPRequest with method/path/version
	return &httpReq, nil
}
