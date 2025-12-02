package io

import (
	"cli-t/internal/shared/logger"
	"fmt"
	"io"
	"os"
)

func WriteOutput(data []byte, inputPath, outputPath string, stdout io.Writer, extension string) error {
	logger.Debug("writing output", "input_size", len(data))

	// Determine output
	if outputPath == "" {
		if inputPath == "stdin" {
			// Write to stdout
			logger.Verbose("Writing output data to stdout")
			data = append(data, '\n')
			_, err := stdout.Write(data)
			return err
		}
		// Auto-generate filename
		outputPath = inputPath + fmt.Sprintf(".%s", extension)
	}

	// Write to file
	logger.Verbose("Writing output file", "path", outputPath)
	if err := os.WriteFile(outputPath, data, 0644); err != nil {
		return fmt.Errorf("failed to write output: %w", err)
	}

	return nil
}

type Writer io.Writer

func GetOutputWriter(stdout io.Writer, outputPath string) (io.Writer, func(), error) {
	// No output file specified, use stdout
	if outputPath == "" {
		return stdout, func() {}, nil
	}

	// Create output file
	outFile, err := os.Create(outputPath)
	if err != nil {
		return nil, func() {}, err
	}

	// Return cleanup function, don't defer here!
	cleanup := func() { outFile.Close() }
	return outFile, cleanup, nil
}
