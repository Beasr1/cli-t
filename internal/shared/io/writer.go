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
			_, err := stdout.Write(data)
			return err
		}
		// Auto-generate filename
		outputPath = inputPath + ".sort"
	}

	// Write to file
	logger.Verbose("Writing output file", "path", outputPath)
	if err := os.WriteFile(outputPath, data, 0644); err != nil {
		return fmt.Errorf("failed to write output: %w", err)
	}

	return nil
}

type Writer io.Writer
