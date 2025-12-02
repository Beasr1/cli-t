package io

import (
	"cli-t/internal/command"
	"cli-t/internal/shared/file"
	"cli-t/internal/shared/logger"
	"fmt"
	"io"
	"os"
)

type Reader io.Reader

// readInput reads from stdin or file, returning path and data
// reads all the content at once
func ReadInput(args *command.Args) (path string, data []byte, err error) {
	if len(args.Positional) == 0 {
		// Read from stdin
		logger.Verbose("Reading from stdin")
		data, err = io.ReadAll(args.Stdin)
		if err != nil {
			return "", nil, fmt.Errorf("error reading from stdin: %w", err)
		}
		return "stdin", data, nil
	}

	// Read from file
	path = args.Positional[0]
	logger.Verbose("Reading from file", "path", path)

	data, err = file.ReadBytes(path)
	if err != nil {
		return "", nil, fmt.Errorf("error reading file: %w", err)
	}

	return path, data, nil
}

func GetInputReader(stdin io.Reader, inputPath string) (io.Reader, func(), error) {
	if inputPath == "" {
		return stdin, func() {}, nil // No cleanup needed
	}

	file, err := os.Open(inputPath)
	if err != nil {
		return nil, nil, err
	}

	// Return cleanup function, don't defer here!
	cleanup := func() { file.Close() }
	return file, cleanup, nil
}
