package protocol

import (
	"fmt"
)

func (s SimpleString) Serialize() []byte {
	return []byte(fmt.Sprintf("+%s\r\n", s.Value))
}

func (err Error) Serialize() []byte {
	return []byte(fmt.Sprintf("-%s\r\n", err.Message))
}

func (i Integer) Serialize() []byte {
	return []byte(fmt.Sprintf(":%d\r\n", i.Value))
}
func (b BulkString) Serialize() []byte {
	// Check if null
	if b.IsNull {
		return []byte("$-1\r\n")
	}

	// Normal case: $<length>\r\n<data>\r\n
	return []byte(fmt.Sprintf("$%d\r\n%s\r\n", len(b.Value), b.Value))
}

func (a Array) Serialize() []byte {
	// Check if null
	if a.IsNull {
		return []byte("*-1\r\n")
	}

	// Start with count: *<count>\r\n
	// appending string or byte is one and same thing
	result := []byte(fmt.Sprintf("*%d\r\n", len(a.Elements)))

	for _, elem := range a.Elements {
		result = append(result, elem.Serialize()...)
	}

	return result
}
