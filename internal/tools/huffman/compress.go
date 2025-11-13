package huffman

import (
	"bytes"
	"cli-t/internal/tools/huffman/algorithm"
	"encoding/binary"
	"fmt"
)

const MAGIC_STRING = "HUFF_ATMKBFG"

// CompressedHeader stores metadata for decompression
type CompressedHeader struct {
	OriginalSize  uint32          // Original file size
	Padding       uint8           // Number of padding bits
	FreqTableSize uint16          // Number of unique characters
	FreqTable     map[rune]uint32 // Character frequencies
}

func Compress(content string) ([]byte, error) {
	if len(content) == 0 {
		return nil, fmt.Errorf("empty content")
	}

	// Build frequency table
	frequencyTable := algorithm.BuildFrequencyTable(content)

	// Build huffman tree
	huffmanTree := algorithm.BuildHuffmanTree(frequencyTable)
	if huffmanTree == nil {
		return nil, fmt.Errorf("failed to build huffman tree")
	}

	// Generate codes
	codeMap := huffmanTree.GenerateCodes()

	// Encode data
	encodedString := algorithm.EncodeData(content, codeMap)

	// Pack bits
	packedData, padding := algorithm.PackBits(encodedString)

	// Create header
	header := CompressedHeader{
		OriginalSize:  uint32(len(content)),
		Padding:       uint8(padding),
		FreqTableSize: uint16(len(frequencyTable)),
		FreqTable:     make(map[rune]uint32),
	}

	// Convert frequencies to uint32 for serialization
	for char, freq := range frequencyTable {
		header.FreqTable[char] = uint32(freq)
	}

	// Serialize header + data
	return serializeCompressed(header, packedData)
}

func serializeCompressed(header CompressedHeader, data []byte) ([]byte, error) {
	buf := new(bytes.Buffer)

	// Write magic number (for file identification)
	// for "HUFF_ATMKBFG" = 12 bytes
	magic := []byte(MAGIC_STRING)
	buf.Write(magic)

	// Write original size (1 int32) = 4 bytes
	binary.Write(buf, binary.LittleEndian, header.OriginalSize)

	// Write padding 0 to 7 = 1 byte
	buf.WriteByte(header.Padding)

	// Write frequency table size
	// can list all characters for simplicity but lets have flexiblity (0 to all characters) = 2 bytes
	// uint8 might miss some edge case that is unicode support (thousands)
	// 1 byte is not worth the risk of breaking on edge cases
	binary.Write(buf, binary.LittleEndian, header.FreqTableSize)

	// Write frequency table
	for char, freq := range header.FreqTable {
		// Write character (as int32)
		binary.Write(buf, binary.LittleEndian, int32(char))
		// Write frequency
		binary.Write(buf, binary.LittleEndian, freq)
	}

	// Write compressed data
	buf.Write(data)

	return buf.Bytes(), nil
}
