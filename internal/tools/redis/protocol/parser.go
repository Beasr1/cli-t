package protocol

import (
	"bytes"
	"errors"
	"fmt"
	"strconv"
)

// TODO: much later Parser Allocations ans use buffer
func Parse(data []byte) (RESPValue, int, error) {
	if len(data) == 0 {
		return nil, 0, errors.New("empty input")
	}

	switch data[0] {
	case '+':
		return parseSimpleString(data)
	case '-':
		return parseError(data)
	case ':':
		return parseInteger(data)
	case '$':
		return parseBulkString(data)
	case '*':
		return parseArray(data)
	default:
		return nil, 0, fmt.Errorf("unknown type byte: %c", data[0])
	}
}

func parseSimpleString(data []byte) (SimpleString, int, error) {
	// 1. Find where \r\n is
	idx := bytes.Index(data, []byte("\r\n"))
	if idx == -1 {
		return SimpleString{}, 0, errors.New("missing \\r\\n terminator")
	}

	// 2. Check it starts with '+'
	if data[0] != '+' {
		return SimpleString{}, 0, errors.New("not a simple string")
	}

	// 3. Extract string between '+' and '\r\n'
	value := string(data[1:idx])
	consumed := idx + 2

	return SimpleString{Value: value}, consumed, nil
}

func parseError(data []byte) (Error, int, error) {
	// 1. Find where \r\n is
	idx := bytes.Index(data, []byte("\r\n"))
	if idx == -1 {
		return Error{}, 0, errors.New("missing \\r\\n terminator")
	}

	// 2. Check it starts with '-'
	if data[0] != '-' {
		return Error{}, 0, errors.New("not a error")
	}

	// 3. Extract string between '-' and '\r\n'
	message := string(data[1:idx])
	consumed := idx + 2

	return Error{Message: message}, consumed, nil
}

func parseInteger(data []byte) (Integer, int, error) {
	// 1. Find where \r\n is
	idx := bytes.Index(data, []byte("\r\n"))
	if idx == -1 {
		return Integer{}, 0, errors.New("missing \\r\\n terminator")
	}

	// 2. Check it starts with ':'
	if data[0] != ':' {
		return Integer{}, 0, errors.New("not a integer")
	}

	// 3. Extract string between ':' and '\r\n'
	value, err := strconv.ParseInt(string(data[1:idx]), 10, 64)
	if err != nil {
		return Integer{}, 0, errors.New("not a integer")
	}

	consumed := idx + 2
	return Integer{Value: value}, consumed, nil
}

func parseBulkString(data []byte) (BulkString, int, error) {
	// 1. Check starts with '$'
	if data[0] != '$' {
		return BulkString{}, 0, errors.New("not a bulk string")
	}

	// 2. Find first \r\n (end of length line)
	idx := bytes.Index(data, []byte("\r\n"))
	if idx == -1 {
		return BulkString{}, 0, errors.New("missing \\r\\n terminator")
	}

	// 3. Parse the length
	// bytes.Split(data, []byte("\r\n")) : will not work since bulk string can contain \r\n inside it
	lengthStr := string(data[1:idx]) // Everything between $ and \r\n
	length, err := strconv.Atoi(lengthStr)
	if err != nil {
		return BulkString{}, 0, fmt.Errorf("invalid length: %v", err)
	}

	// 4. Handle null case
	if length == -1 {
		consumed := idx + 2
		return BulkString{IsNull: true}, consumed, nil
	}

	// 5. Extract the actual data (after first \r\n)
	dataStart := idx + 2 // Skip past first \r\n
	dataEnd := dataStart + length

	// 6. Validate we have enough data
	if dataEnd > len(data) {
		return BulkString{}, 0, errors.New("data shorter than specified length")
	}

	value := string(data[dataStart:dataEnd])

	// 7. Check for final \r\n
	if dataEnd+2 > len(data) || string(data[dataEnd:dataEnd+2]) != "\r\n" {
		return BulkString{}, 0, errors.New("missing final \\r\\n")
	}

	consumed := dataEnd + 2
	return BulkString{Value: value}, consumed, nil
}

func parseArray(data []byte) (Array, int, error) {
	// 1. Check starts with '*'
	if data[0] != '*' {
		return Array{}, 0, errors.New("not a array")
	}

	// 2. Find first \r\n (end of length line)
	idx := bytes.Index(data, []byte("\r\n"))
	if idx == -1 {
		return Array{}, 0, errors.New("missing \\r\\n terminator")
	}

	// 3. Parse the length
	lengthStr := string(data[1:idx]) // Everything between $ and \r\n
	length, err := strconv.Atoi(lengthStr)
	if err != nil {
		return Array{}, 0, fmt.Errorf("invalid length: %v", err)
	}

	// 4. Handle null case
	if length == -1 {
		consumed := idx + 2
		return Array{IsNull: true}, consumed, nil
	}

	elements := make([]RESPValue, 0, length)
	pos := idx + 2 // Start after "*<count>\r\n"

	for i := 0; i < length; i++ {
		elem, consumed, err := Parse(data[pos:]) // Recursive!
		if err != nil {
			return Array{}, 0, err
		}
		elements = append(elements, elem)
		pos += consumed
	}

	totalConsumed := pos // All bytes from start

	return Array{Elements: elements}, totalConsumed, nil
}
