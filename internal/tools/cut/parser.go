package cut

import (
	"fmt"
	"strconv"
	"strings"
)

type Selection struct {
	Start int // -1 means "from beginning"
	End   int // -1 means "to end"
}

// ParseList parses a list specification like "1", "1-3", "1,3,5", "-3", "5-"
func ParseList(input string) ([]Selection, error) {
	if input == "" {
		return nil, fmt.Errorf("empty list specification")
	}

	var selections []Selection

	// Split by comma
	parts := strings.Split(input, ",")

	for _, part := range parts {
		part = strings.TrimSpace(part)
		if part == "" {
			continue
		}

		// Check if it's a range (contains '-')
		if strings.Contains(part, "-") {
			sel, err := parseRange(part)
			if err != nil {
				return nil, err
			}
			selections = append(selections, sel)
		} else {
			// Single number
			num, err := strconv.Atoi(part)
			if err != nil {
				return nil, fmt.Errorf("invalid number: %s", part)
			}
			if num < 1 {
				return nil, fmt.Errorf("numbers must be >= 1, got %d", num)
			}
			selections = append(selections, Selection{Start: num, End: num})
		}
	}

	if len(selections) == 0 {
		return nil, fmt.Errorf("no valid selections found")
	}

	return selections, nil
}

// parseRange parses a range like "1-3", "-3", "5-"
func parseRange(s string) (Selection, error) {
	// Handle open-ended ranges
	if strings.HasPrefix(s, "-") {
		// "-3" means from beginning to 3
		endStr := strings.TrimPrefix(s, "-")
		end, err := strconv.Atoi(endStr)
		if err != nil {
			return Selection{}, fmt.Errorf("invalid range end: %s", endStr)
		}
		if end < 1 {
			return Selection{}, fmt.Errorf("range end must be >= 1, got %d", end)
		}
		return Selection{Start: -1, End: end}, nil
	}

	if strings.HasSuffix(s, "-") {
		// "5-" means from 5 to end
		startStr := strings.TrimSuffix(s, "-")
		start, err := strconv.Atoi(startStr)
		if err != nil {
			return Selection{}, fmt.Errorf("invalid range start: %s", startStr)
		}
		if start < 1 {
			return Selection{}, fmt.Errorf("range start must be >= 1, got %d", start)
		}
		return Selection{Start: start, End: -1}, nil
	}

	// Normal range "1-3"
	parts := strings.Split(s, "-")
	if len(parts) != 2 {
		return Selection{}, fmt.Errorf("invalid range format: %s", s)
	}

	start, err := strconv.Atoi(strings.TrimSpace(parts[0]))
	if err != nil {
		return Selection{}, fmt.Errorf("invalid range start: %s", parts[0])
	}

	end, err := strconv.Atoi(strings.TrimSpace(parts[1]))
	if err != nil {
		return Selection{}, fmt.Errorf("invalid range end: %s", parts[1])
	}

	if start < 1 || end < 1 {
		return Selection{}, fmt.Errorf("range values must be >= 1")
	}

	if start > end {
		return Selection{}, fmt.Errorf("range start (%d) cannot be greater than end (%d)", start, end)
	}

	return Selection{Start: start, End: end}, nil
}
