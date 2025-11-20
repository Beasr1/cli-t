package protocol

// RESP represents a Redis Serialization Protocol value
type RESPValue interface {
	Serialize() []byte
}

// Simple String: "+OK\r\n"
type SimpleString struct {
	Value string
}

// Error: "-ERR unknown command\r\n"
type Error struct {
	Message string
}

// Integer: ":42\r\n"
type Integer struct {
	Value int64
}

// Bulk String: "$5\r\nhello\r\n" or "$-1\r\n" (null)
type BulkString struct {
	Value  string
	IsNull bool
}

// Array: "*2\r\n$3\r\nfoo\r\n$3\r\nbar\r\n" or "*-1\r\n" (null)
type Array struct {
	IsNull   bool
	Elements []RESPValue
}

// Examples of what we need to represent:
// "+OK\r\n"        → SimpleString{value: "OK"}
// ":42\r\n"        → Integer{value: 42}
// "$-1\r\n"        → ??? (null bulk string)
// "*2\r\n$3\r\nfoo\r\n$3\r\nbar\r\n" → Array{[BulkString, BulkString]}
