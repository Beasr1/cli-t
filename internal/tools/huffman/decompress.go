package huffman

import (
	"bytes"
	"cli-t/internal/tools/huffman/algorithm"
	"encoding/binary"
	"fmt"
	"strings"
)

func Decompress(compressedData []byte) (string, error) {
	if len(compressedData) < 4 {
		return "", fmt.Errorf("invalid compressed data")
	}

	// Deserialize
	header, data, err := deserializeCompressed(compressedData)
	if err != nil {
		return "", err
	}

	// Rebuild frequency table
	freqTable := make(map[rune]int)
	for char, freq := range header.FreqTable {
		freqTable[char] = int(freq)
	}

	// Rebuild huffman tree
	huffmanTree := algorithm.BuildHuffmanTree(freqTable)
	if huffmanTree == nil {
		return "", fmt.Errorf("failed to rebuild huffman tree")
	}

	// Unpack bits
	bitString := unpackBits(data, header.Padding)

	// Decode data
	decoded := huffmanTree.Decode(bitString)

	return decoded, nil
}

func deserializeCompressed(data []byte) (CompressedHeader, []byte, error) {
	buf := bytes.NewReader(data)
	header := CompressedHeader{}

	// Read magic number
	magic := make([]byte, len(MAGIC_STRING))
	buf.Read(magic)
	if string(magic) != MAGIC_STRING {
		return header, nil, fmt.Errorf("invalid file format")
	}

	// Read original size
	binary.Read(buf, binary.LittleEndian, &header.OriginalSize)

	// Read padding
	padding, _ := buf.ReadByte()
	header.Padding = padding

	// Read frequency table size
	binary.Read(buf, binary.LittleEndian, &header.FreqTableSize)

	// Read frequency table
	header.FreqTable = make(map[rune]uint32)
	for i := 0; i < int(header.FreqTableSize); i++ {
		var char int32
		var freq uint32
		binary.Read(buf, binary.LittleEndian, &char)
		binary.Read(buf, binary.LittleEndian, &freq)
		header.FreqTable[rune(char)] = freq
	}

	// Read remaining data (compressed content)
	remaining := make([]byte, buf.Len())
	buf.Read(remaining)

	return header, remaining, nil
}

func unpackBits(data []byte, padding uint8) string {
	var bitString strings.Builder

	for _, b := range data {
		// Convert byte to 8-bit string
		for bit := 7; bit >= 0; bit-- {
			if b&(1<<bit) != 0 {
				bitString.WriteByte('1')
			} else {
				bitString.WriteByte('0')
			}
		}
	}

	// Remove padding from last byte
	bits := bitString.String()
	if padding > 0 && len(bits) >= int(padding) {
		bits = bits[:len(bits)-int(padding)]
	}

	return bits
}
