package huffman

import (
	"container/heap"
	"fmt"
	"strconv"
	"strings"
)

// will need to create heap of HuffmanNodes then make a tree out of it

/*
BuildFrequencyTable counts the occurrence of each rune in the input string
Returns a map where keys are runes and values are their frequencies

NOTE: compression on []byte is more efficent for large files (still focusing on runes first)
*/
func BuildFrequencyTable(content string) map[rune]int {
	if len(content) == 0 {
		return make(map[rune]int) // or return nil, depending on your design
	}

	// NOTE: can preallocate a map with rough estimates (optimistic approach)
	frequency := map[rune]int{}
	for _, run := range content {
		frequency[run]++
	}

	return frequency
}

/*
BuildHuffmanTree constructs a Huffman tree from character frequencies.
It uses a min-heap to repeatedly merge the two nodes with lowest frequencies
until a single tree remains.

Parameters:
  - frequencies: map of characters to their occurrence counts

Returns:
  - *HuffmanTree: the root of the constructed tree, or nil if frequencies is empty

Algorithm:
 1. Create leaf nodes for each character and add to min-heap
 2. While heap has more than 1 node:
    - Extract two nodes with minimum frequency
    - Create new internal node with sum of frequencies
    - Add internal node back to heap
 3. The last remaining node is the root
*/
func BuildHuffmanTree(frequencies map[rune]int) *HuffmanTree {
	minHeap := HuffmanMinHeap{}
	for key, val := range frequencies {
		hn := NewHuffmanNode(key, val, true) // returns points of the actual value

		// will internally call minHeap.Push() then heapifies
		heap.Push(&minHeap, hn)
	}

	// Build tree
	for minHeap.Len() > 1 {
		left := heap.Pop(&minHeap).(*HuffmanNode)
		right := heap.Pop(&minHeap).(*HuffmanNode)

		// 0 value stands for no-character
		merged := NewHuffmanNode(0, left.frequency+right.frequency, false)
		merged.left = left
		merged.right = right
		heap.Push(&minHeap, merged)
	}

	// Get root
	if minHeap.Len() == 1 {
		return &HuffmanTree{rootNode: heap.Pop(&minHeap).(*HuffmanNode)}
	}

	return nil
}

func EncodeData(content string, codeMap map[rune]string) string {
	var encodedData strings.Builder
	for _, run := range content {
		encodedData.WriteString(codeMap[run])
	}

	// 	If we store as string: "000101011" = 9 bytes (72 bits!)
	// If we pack into bytes: [00010101] [1_______] = 2 bytes (16 bits)
	encodedString := encodedData.String()

	return encodedString
}

// PackBits converts a string of '0' and '1' characters into packed bytes.
// It processes the bit string in 8-bit chunks, converting each chunk to a byte.
// If the last chunk has fewer than 8 bits, it pads with zeros on the right.
//
// Parameters:
//   - bitString: string containing only '0' and '1' characters
//
// Returns:
//   - packed: slice of bytes containing the packed binary data
//   - padding: number of padding bits added to the last byte (0-7)
//
// Example:
//
//	PackBits("0101101111") returns ([0x5B, 0xE0], 5)
//	where 0x5B = 01011011 and 0xE0 = 11100000 (with 5 padding zeros)
func PackBits(bitString string) (packed []byte, padding int) {
	// Pre-allocate slice for better performance
	packed = make([]byte, 0, (len(bitString)+7)/8)

	// Process in 8-bit chunks
	for i := 0; i < len(bitString); i += 8 {
		end := i + 8
		if end > len(bitString) {
			end = len(bitString)
		}

		chunk := bitString[i:end]

		// Pad the last chunk if needed
		if len(chunk) < 8 {
			padding = 8 - len(chunk)
			chunk = chunk + strings.Repeat("0", padding)
		}

		// Convert binary string to byte
		byteVal, err := strconv.ParseUint(chunk, 2, 8)
		if err != nil {
			// This should never happen with valid input
			panic(fmt.Sprintf("invalid bit string at position %d: %v", i, err))
		}

		packed = append(packed, byte(byteVal))
	}

	return packed, padding
}
