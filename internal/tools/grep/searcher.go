package grep

import (
	"bufio"
	"cli-t/internal/shared/logger"
	"fmt"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"regexp"
)

// SearchFile searches a single file line-by-line
// Returns: foundMatch (bool), error
func SearchFile(filepath, pattern string, showFileName, invert bool) (bool, error) {
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
		return false, err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)

	// Compile ONCE (outside the loop)
	re, err := regexp.Compile(pattern)
	if err != nil {
		logger.Error("error", "error", err)
		return false, err
	}

	prefix := ""
	if showFileName {
		prefix = filepath + ":"
	}

	// optionally, resize scanner's capacity for lines over 64K, see next example
	m := false
	for scanner.Scan() {
		line := scanner.Text()
		shouldPrint := re.MatchString(line) != invert // or use XOR
		if shouldPrint {
			m = true
			fmt.Println(prefix + line)
		}
	}

	if err := scanner.Err(); err != nil {
		return false, err
	}

	return m, nil
}

func SearchDirectory(dir string, pattern string, invert bool) (bool, error) {
	// Use filepath.WalkDir to visit each file
	// For each file (not directory), call SearchFile
	matched := false
	err := filepath.WalkDir(dir, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if d.IsDir() {
			return nil // Skip, we only want files
		}

		match, err := SearchFile(path, pattern, true, invert)
		if err != nil {
			return err
		}
		matched = match || matched
		return nil
	})

	return matched, err
}

func SearchStdin(reader io.Reader, pattern string, invert bool) (bool, error) {
	scanner := bufio.NewScanner(reader)
	re, err := regexp.Compile(pattern)
	if err != nil {
		logger.Error("error", "error", err)
		return false, err
	}

	prefix := ""

	// optionally, resize scanner's capacity for lines over 64K, see next example
	m := false
	for scanner.Scan() {
		line := scanner.Text()
		shouldPrint := re.MatchString(line) != invert // or use XOR
		if shouldPrint {
			m = true
			fmt.Println(prefix + line)
		}
	}

	if err := scanner.Err(); err != nil {
		return false, err
	}

	return m, nil
}
