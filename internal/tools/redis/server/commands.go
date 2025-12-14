package server

import (
	"cli-t/internal/shared/logger"
	inmemory "cli-t/internal/shared/store/inmemory"
	"cli-t/internal/tools/redis/protocol"

	"strconv"
	"strings"
	"time"
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
	case "TTL":
		return s.handleTtl(arr.Elements)
	case "EXPIRE":
		return s.handleExpire(arr.Elements)
	case "EXISTS":
		return s.handleExists(arr.Elements)
	case "DEL":
		return s.handleDelete(arr.Elements)
	case "INCR":
		return s.handleIncr(arr.Elements)
	case "DECR":
		return s.handleDecr(arr.Elements)
	case "LPUSH":
		return s.handleLPush(arr.Elements)
	case "RPUSH":
		return s.handleRPush(arr.Elements)
	case "LRANGE":
		return s.handleLRange(arr.Elements)
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
	if len(args) < 3 {
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

	// Parse optional EX/PX flags
	var expiresAt *time.Time
	for i := 3; i < len(args); i++ {
		flag, ok := args[i].(protocol.BulkString)
		if !ok {
			return protocol.Error{Message: "ERR syntax error"}
		}
		f := strings.ToUpper(flag.Value)

		switch f {
		case "EX":
			if i+1 >= len(args) {
				return protocol.Error{Message: "ERR syntax error"}
			}

			secondsArg, ok := args[i+1].(protocol.BulkString)
			if !ok {
				return protocol.Error{Message: "ERR value is not an integer"}
			}

			// Parse seconds, create time.Now().Add(seconds)
			seconds, err := strconv.Atoi(secondsArg.Value)
			if err != nil || seconds <= 0 {
				return protocol.Error{Message: "ERR value is not an integer or out of range"}
			}

			expiry := time.Now().Add(time.Duration(seconds) * time.Second)
			expiresAt = &expiry
			i++ // Skip next arg

			logger.Debug("SET with expiry", "key", key.Value, "seconds", seconds)
		case "PX":
			if i+1 >= len(args) {
				return protocol.Error{Message: "ERR syntax error"}
			}

			milliSecondsArg, ok := args[i+1].(protocol.BulkString)
			if !ok {
				return protocol.Error{Message: "ERR value is not an integer"}
			}

			// Create expiry: time.Now().Add(time.Duration(millis) * time.Millisecond)
			milliSeconds, err := strconv.Atoi(milliSecondsArg.Value)
			if err != nil || milliSeconds <= 0 {
				return protocol.Error{Message: "ERR value is not an integer or out of range"}
			}

			expiry := time.Now().Add(time.Duration(milliSeconds) * time.Millisecond)
			expiresAt = &expiry
			i++ // Skip next arg

			logger.Debug("SET with PX", "key", key.Value, "milliSeconds", milliSeconds)
		case "EXAT":
			if i+1 >= len(args) {
				return protocol.Error{Message: "ERR syntax error"}
			}

			timestampArg, ok := args[i+1].(protocol.BulkString)
			if !ok {
				return protocol.Error{Message: "ERR value is not an integer"}
			}

			timestampInt, err := strconv.ParseInt(timestampArg.Value, 10, 64)
			if err != nil {
				return protocol.Error{Message: "ERR value is not an integer"}
			}

			timestamp := time.Unix(timestampInt, 0)
			if time.Until(timestamp) <= 0 {
				return protocol.Error{Message: "ERR invalid expire time in 'set' command"}
			}

			expiresAt = &timestamp
			i++ // Skip next arg

			logger.Debug("SET with EXAT", "key", key.Value, "timestamp", timestamp)
		case "PXAT":
			if i+1 >= len(args) {
				return protocol.Error{Message: "ERR syntax error"}
			}

			timestampMillisArg, ok := args[i+1].(protocol.BulkString)
			if !ok {
				return protocol.Error{Message: "ERR value is not an integer"}
			}

			timestampMs, err := strconv.ParseInt(timestampMillisArg.Value, 10, 64)
			if err != nil {
				return protocol.Error{Message: "ERR value is not an integer"}
			}

			timestamp := time.Unix(timestampMs/1000, (timestampMs%1000)*1000000)
			if time.Until(timestamp) <= 0 {
				return protocol.Error{Message: "ERR invalid expire time in 'set' command"}
			}

			expiresAt = &timestamp
			i++ // Skip next arg

			logger.Debug("SET with PXAT", "key", key.Value, "timestamp", timestamp)
		default:
			return protocol.Error{Message: "ERR syntax error"}
		}
	}

	s.store.Set(key.Value, inmemory.StoreValue{
		Data:      value.Value,
		ExpiresAt: expiresAt,
		Type:      inmemory.TypeString,
	})

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
		Value:  value.Data,
		IsNull: false, //  Not null when found
	}
}

func (s *Server) handleTtl(args []protocol.RESPValue) protocol.RESPValue {
	if len(args) != 2 { //  GET key (2 args total)
		return protocol.Error{Message: "ERR wrong number of arguments for 'ttl' command"}
	}

	// Safe type assertion
	key, ok := args[1].(protocol.BulkString)
	if !ok {
		return protocol.Error{Message: "ERR key must be a string"}
	}

	ttl := s.store.GetTTL(key.Value)
	return protocol.Integer{
		Value: ttl,
	}

}

func (s *Server) handleExpire(args []protocol.RESPValue) protocol.RESPValue {
	if len(args) != 3 {
		return protocol.Error{Message: "ERR wrong number of arguments for 'expire' command"}
	}

	// Safe type assertions
	key, ok := args[1].(protocol.BulkString)
	if !ok {
		return protocol.Error{Message: "ERR key must be a string"}
	}

	expireBulk, ok := args[2].(protocol.BulkString)
	if !ok {
		return protocol.Error{Message: "ERR value must be a string"}
	}

	expire, err := strconv.Atoi(expireBulk.Value)
	if err != nil || expire < 0 {
		return protocol.Error{Message: "ERR value is not an integer or out of range"}
	}

	success := s.store.SetExpiry(key.Value, expire)
	if success {
		return protocol.Integer{Value: 1}
	} else {
		return protocol.Integer{Value: 0}
	}

}

func (s *Server) handleExists(args []protocol.RESPValue) protocol.RESPValue {
	if len(args) < 2 {
		return protocol.Error{Message: "ERR wrong number of arguments for 'EXISTS' command"}
	}

	keys := []string{}

	for i := 1; i < len(args); i++ {
		// Safe type assertions
		key, ok := args[i].(protocol.BulkString)
		if !ok {
			return protocol.Error{Message: "ERR key must be a string"}
		}
		keys = append(keys, key.Value)
	}

	exists := s.store.Exists(keys...)
	return protocol.Integer{Value: int64(exists)}
}

func (s *Server) handleDelete(args []protocol.RESPValue) protocol.RESPValue {
	if len(args) < 2 {
		return protocol.Error{Message: "ERR wrong number of arguments for 'DEL' command"}
	}

	keys := []string{}

	for i := 1; i < len(args); i++ {
		// Safe type assertions
		key, ok := args[i].(protocol.BulkString)
		if !ok {
			return protocol.Error{Message: "ERR key must be a string"}
		}
		keys = append(keys, key.Value)
	}

	count := s.store.Delete(keys...)
	return protocol.Integer{Value: int64(count)}
}

func (s *Server) handleIncr(args []protocol.RESPValue) protocol.RESPValue {
	// Check arg count
	// Extract key
	// Call s.store.Incr(key)
	// If error, return protocol.Error
	// Otherwise, return protocol.Integer with new value
	if len(args) != 2 { //  GET key (2 args total)
		return protocol.Error{Message: "ERR wrong number of arguments for 'incr' command"}
	}

	// Safe type assertion
	key, ok := args[1].(protocol.BulkString)
	if !ok {
		return protocol.Error{Message: "ERR key must be a string"}
	}

	newVal, err := s.store.Incr(key.Value)
	if err != nil {
		return protocol.Error{Message: err.Error()}
	}

	return protocol.Integer{
		Value: newVal,
	}
}

func (s *Server) handleDecr(args []protocol.RESPValue) protocol.RESPValue {
	// Check arg count
	// Extract key
	// Call s.store.Decr(key)
	// If error, return protocol.Error
	// Otherwise, return protocol.Integer with new value
	if len(args) != 2 { //  GET key (2 args total)
		return protocol.Error{Message: "ERR wrong number of arguments for 'decr' command"}
	}

	// Safe type assertion
	key, ok := args[1].(protocol.BulkString)
	if !ok {
		return protocol.Error{Message: "ERR key must be a string"}
	}

	newVal, err := s.store.Decr(key.Value)
	if err != nil {
		return protocol.Error{Message: err.Error()}
	}

	return protocol.Integer{
		Value: newVal,
	}
}

func (s *Server) handleLPush(args []protocol.RESPValue) protocol.RESPValue {
	if len(args) < 3 { // LPUSH key value1 [value2 ...]
		return protocol.Error{Message: "ERR wrong number of arguments for 'lpush' command"}
	}

	key, ok := args[1].(protocol.BulkString)
	if !ok {
		return protocol.Error{Message: "ERR key must be a string"}
	}

	// Collect all values
	values := []string{}
	for i := 2; i < len(args); i++ {
		val, ok := args[i].(protocol.BulkString)
		if !ok {
			return protocol.Error{Message: "ERR value must be a string"}
		}
		values = append(values, val.Value)
	}

	count, err := s.store.LPush(key.Value, values...)
	if err != nil {
		return protocol.Error{Message: err.Error()}
	}

	return protocol.Integer{Value: count}
}

func (s *Server) handleRPush(args []protocol.RESPValue) protocol.RESPValue {
	if len(args) < 3 { // RPUSH key value1 [value2 ...]
		return protocol.Error{Message: "ERR wrong number of arguments for 'lpush' command"}
	}

	key, ok := args[1].(protocol.BulkString)
	if !ok {
		return protocol.Error{Message: "ERR key must be a string"}
	}

	// Collect all values
	values := []string{}
	for i := 2; i < len(args); i++ {
		val, ok := args[i].(protocol.BulkString)
		if !ok {
			return protocol.Error{Message: "ERR value must be a string"}
		}
		values = append(values, val.Value)
	}

	count, err := s.store.RPush(key.Value, values...)
	if err != nil {
		return protocol.Error{Message: err.Error()}
	}

	return protocol.Integer{Value: count}
}

func (s *Server) handleLRange(args []protocol.RESPValue) protocol.RESPValue {
	if len(args) != 4 { // LRANGE key start stop
		return protocol.Error{Message: "ERR wrong number of arguments for 'lrange' command"}
	}

	key, ok := args[1].(protocol.BulkString)
	if !ok {
		return protocol.Error{Message: "ERR key must be a string"}
	}

	startBulk, ok := args[2].(protocol.BulkString)
	if !ok {
		return protocol.Error{Message: "ERR value must be a string"}
	}

	start, err := strconv.Atoi(startBulk.Value)
	if err != nil {
		return protocol.Error{Message: "ERR value is not an integer or out of range"}
	}

	stopBulk, ok := args[3].(protocol.BulkString)
	if !ok {
		return protocol.Error{Message: "ERR value must be a string"}
	}

	stop, err := strconv.Atoi(stopBulk.Value)
	if err != nil {
		return protocol.Error{Message: "ERR value is not an integer or out of range"}
	}

	// Redis allows negative indice
	result, err := s.store.LRange(key.Value, start, stop)
	if err != nil {
		return protocol.Error{Message: err.Error()}
	}

	// Convert []string to Array of BulkStrings
	elements := make([]protocol.RESPValue, len(result))
	for i, val := range result {
		elements[i] = protocol.BulkString{Value: val}
	}

	return protocol.Array{Elements: elements}
}
