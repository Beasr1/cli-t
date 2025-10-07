package cut

import "strings"

// ExtractFields extracts selected fields from a line based on delimiter
func ExtractFields(line string, delimiter rune, selections []Selection) string {
	// Split line into fields
	fields := strings.Split(line, string(delimiter))
	totalFields := len(fields)

	// Track which fields to extract (avoid duplicates)
	shouldExtract := make(map[int]bool)

	// Process each selection
	for _, sel := range selections {
		start := sel.Start
		end := sel.End

		// Handle -1 for "from beginning"
		if start == -1 {
			start = 1
		}

		// Handle -1 for "to end"
		if end == -1 {
			end = totalFields
		}

		// Mark fields in this range
		for fieldNum := start; fieldNum <= end; fieldNum++ {
			// Only mark if field exists (1-indexed)
			if fieldNum >= 1 && fieldNum <= totalFields {
				shouldExtract[fieldNum] = true
			}
		}
	}

	// Extract marked fields in order
	var result []string
	for fieldNum := 1; fieldNum <= totalFields; fieldNum++ {
		if shouldExtract[fieldNum] {
			// fieldNum is 1-indexed, array is 0-indexed
			result = append(result, fields[fieldNum-1])
		}
	}

	// Join with delimiter
	return strings.Join(result, string(delimiter))
}

// ExtractChars extracts selected characters from a line
func ExtractChars(line string, selections []Selection) string {
	// Convert to runes for proper Unicode handling
	runes := []rune(line)
	totalChars := len(runes)

	// Track which characters to extract
	shouldExtract := make(map[int]bool)

	// Process each selection
	for _, sel := range selections {
		start := sel.Start
		end := sel.End

		// Handle -1 for "from beginning"
		if start == -1 {
			start = 1
		}

		// Handle -1 for "to end"
		if end == -1 {
			end = totalChars
		}

		// Mark characters in this range
		for charNum := start; charNum <= end; charNum++ {
			// Only mark if character exists (1-indexed)
			if charNum >= 1 && charNum <= totalChars {
				shouldExtract[charNum] = true
			}
		}
	}

	// Extract marked characters in order
	var result []rune
	for charNum := 1; charNum <= totalChars; charNum++ {
		if shouldExtract[charNum] {
			// charNum is 1-indexed, array is 0-indexed
			result = append(result, runes[charNum-1])
		}
	}

	return string(result)
}

// ExtractBytes extracts selected bytes from a line
func ExtractBytes(line string, selections []Selection) string {
	// Convert to bytes
	bytes := []byte(line)
	totalBytes := len(bytes)

	// Track which bytes to extract
	shouldExtract := make(map[int]bool)

	// Process each selection
	for _, sel := range selections {
		start := sel.Start
		end := sel.End

		// Handle -1 for "from beginning"
		if start == -1 {
			start = 1
		}

		// Handle -1 for "to end"
		if end == -1 {
			end = totalBytes
		}

		// Mark bytes in this range
		for byteNum := start; byteNum <= end; byteNum++ {
			// Only mark if byte exists (1-indexed)
			if byteNum >= 1 && byteNum <= totalBytes {
				shouldExtract[byteNum] = true
			}
		}
	}

	// Extract marked bytes in order
	var result []byte
	for byteNum := 1; byteNum <= totalBytes; byteNum++ {
		if shouldExtract[byteNum] {
			// byteNum is 1-indexed, array is 0-indexed
			result = append(result, bytes[byteNum-1])
		}
	}

	return string(result)
}
