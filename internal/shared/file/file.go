package file

import (
	"bufio"
	"fmt"
	"io"
	"os"
)

// ReadLines reads a file and returns its lines.
func ReadLines(filePath string) ([]string, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to open file: %w", err)
	}
	defer file.Close()

	var lines []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("failed to scan file: %w", err)
	}

	return lines, nil
}

// ReadContent reads entire file content as string
func ReadContent(filePath string) (string, error) {
	data, err := os.ReadFile(filePath)
	if err != nil {
		return "", fmt.Errorf("failed to read file: %w", err)
	}
	return string(data), nil
}

// ReadBytes reads file content as bytes
func ReadBytes(filePath string) ([]byte, error) {
	return os.ReadFile(filePath)
}

// StreamLines reads a file line by line with a callback
// Good for processing large files without loading all into memory
func StreamLines(filePath string, callback func(line string) error) error {
	file, err := os.Open(filePath)
	if err != nil {
		return fmt.Errorf("failed to open file: %w", err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		if err := callback(scanner.Text()); err != nil {
			return err
		}
	}

	return scanner.Err()
}

// ReadFrom reads from any io.Reader (useful for stdin)
func ReadFrom(r io.Reader) (string, error) {
	data, err := io.ReadAll(r)
	if err != nil {
		return "", err
	}
	return string(data), nil
}
