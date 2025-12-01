package grep

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"regexp"
)

// SearchFile searches a single file line-by-line
// Returns: foundMatch (bool), error
func SearchFile(filepath string, matcher Matcher) (bool, error) {
	// Open file
	// Create bufio.Scanner
	// Loop: scanner.Scan()
	//   - Get line with scanner.Text()
	//   - Check if it matches pattern
	//   - Print if match
	//   - Track if any match found
	// Return whether you found any matches

	file, err := os.Open(filepath)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)

	// optionally, resize scanner's capacity for lines over 64K, see next example
	m := false
	for scanner.Scan() {
		line := scanner.Text()
		matched, err := regexp.Match("", []byte(line))
		if err != nil {
			return false, err
		}
		if !matched {
			continue
		}

		m = true
		fmt.Println(line)
	}

	return m, nil
}
