package server

import (
	"cli-t/internal/tools/redis/protocol"
	"strings"
)

// RESP Protocol Rule: Commands Are Always Arrays of Bulk Strings
func (s *Server) handleCommand(msg protocol.RESPValue) protocol.RESPValue {
	// Redis commands come as Arrays of BulkStrings
	// Example: ["PING"] or ["SET", "key", "value"]

	arr, ok := msg.(protocol.Array)
	if !ok || arr.IsNull || len(arr.Elements) == 0 {
		return protocol.Error{Message: "ERR invalid command format"}
	}

	// First element is the command name
	cmdName, ok := arr.Elements[0].(protocol.BulkString)
	if !ok {
		return protocol.Error{Message: "ERR command must be a bulk string"}
	}

	// Commands are case-insensitive
	cmd := strings.ToUpper(cmdName.Value)

	// Route to appropriate handler
	switch cmd {
	case "PING":
		return s.handlePing(arr.Elements)
	case "ECHO":
		return s.handleEcho(arr.Elements)
	case "SET":
		return s.handleSet(arr.Elements)
	case "GET":
		return s.handleGet(arr.Elements)
	default:
		return protocol.Error{Message: "ERR unknown command '" + cmd + "'"}
	}
}

// Implement these handlers:
func (s *Server) handlePing(args []protocol.RESPValue) protocol.RESPValue {
	return protocol.SimpleString{
		Value: "PONG",
	}
}

func (s *Server) handleEcho(args []protocol.RESPValue) protocol.RESPValue {
	if len(args) != 2 { // ECHO takes exactly 1 argument
		return protocol.Error{Message: "ERR wrong number of arguments for 'echo' command"}
	}

	// Type assert safely
	arg, ok := args[1].(protocol.BulkString)
	if !ok {
		return protocol.Error{Message: "ERR argument must be a string"}
	}

	// Return the argument as BulkString
	return protocol.BulkString{
		Value:  arg.Value,
		IsNull: false,
	}
}

func (s *Server) handleSet(args []protocol.RESPValue) protocol.RESPValue {
	if len(args) != 3 {
		return protocol.Error{Message: "ERR wrong number of arguments for 'set' command"}
	}

	// Safe type assertions
	key, ok := args[1].(protocol.BulkString)
	if !ok {
		return protocol.Error{Message: "ERR key must be a string"}
	}

	value, ok := args[2].(protocol.BulkString)
	if !ok {
		return protocol.Error{Message: "ERR value must be a string"}
	}

	s.store.Set(key.Value, value.Value)
	return protocol.SimpleString{Value: "OK"}
}

func (s *Server) handleGet(args []protocol.RESPValue) protocol.RESPValue {
	if len(args) != 2 { //  GET key (2 args total)
		return protocol.Error{Message: "ERR wrong number of arguments for 'get' command"}
	}

	// Safe type assertion
	key, ok := args[1].(protocol.BulkString)
	if !ok {
		return protocol.Error{Message: "ERR key must be a string"}
	}

	value, exists := s.store.Get(key.Value)
	if !exists {
		// Key not found → return null bulk string
		return protocol.BulkString{
			IsNull: true, //  Null when not found
		}
	}

	// Key found → return the value
	return protocol.BulkString{
		Value:  value,
		IsNull: false, //  Not null when found
	}
}
