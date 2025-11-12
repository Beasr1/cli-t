package huffman

import (
	"cli-t/internal/command"
	"cli-t/internal/shared/io"
	"cli-t/internal/shared/logger"
	"context"
	"fmt"
	"os"
	"strings"
)

type Command struct{}

func New() command.Command {
	return &Command{}
}

func (c *Command) Name() string {
	return "huffman"
}

func (c *Command) Usage() string {
	return "[-d] [-o output] [-k] [file]"
}

func (c *Command) Description() string {
	return "Compress and decompress files using Huffman coding"
}

func (c *Command) DefineFlags() []command.Flag {
	return []command.Flag{
		{
			Name:      "decompress",
			Shorthand: "d",
			Usage:     "Decompress the file",
			Type:      "bool",
			Default:   false,
		},
		{
			Name:      "output",
			Shorthand: "o",
			Usage:     "Output file path (default: stdout or auto-generated)",
			Type:      "string",
			Default:   "",
		},
		{
			Name:      "keep",
			Shorthand: "k",
			Usage:     "Keep (don't delete) input file",
			Type:      "bool",
			Default:   true,
		},
	}
}

func (c *Command) ValidateArgs(args []string) error {
	if len(args) > 1 {
		return fmt.Errorf("too many arguments, expected 0 or 1 file")
	}
	return nil
}

func (c *Command) Execute(ctx context.Context, args *command.Args) error {
	// Parse flags
	decompress, outputPath, keep, err := c.parseFlags(args.Flags)
	if err != nil {
		return err
	}

	logger.Debug("Huffman configuration",
		"decompress", decompress,
		"output", outputPath,
		"keep", keep,
	)

	// Read input
	inputPath, data, err := io.ReadInput(args)
	if err != nil {
		return err
	}

	// Process
	if decompress {
		return c.doDecompress(data, inputPath, outputPath, args.Stdout)
	}
	return c.doCompress(data, inputPath, outputPath, args.Stdout)
}

// parseFlags extracts and validates flag values
func (c *Command) parseFlags(flags map[string]interface{}) (decompress bool, outputPath string, keep bool, err error) {
	decompress, _ = flags["decompress"].(bool)
	outputPath, _ = flags["output"].(string)
	keep, _ = flags["keep"].(bool)

	// Optional: Add validations
	// For example, check if output path is writable
	if outputPath != "" {
		// Check if we can write to this path
		dir := outputPath
		if idx := strings.LastIndex(outputPath, "/"); idx != -1 {
			dir = outputPath[:idx]
		} else {
			dir = "."
		}

		// Check directory exists
		if info, statErr := os.Stat(dir); statErr != nil || !info.IsDir() {
			err = fmt.Errorf("output directory does not exist: %s", dir)
			return
		}
	}

	return
}

// doCompress performs compression
func (c *Command) doCompress(data []byte, inputPath, outputPath string, stdout io.Writer) error {
	logger.Debug("Starting compression", "input_size", len(data))

	// Compress
	compressed, err := Compress(string(data))
	if err != nil {
		return fmt.Errorf("compression failed: %w", err)
	}

	// Determine output
	if outputPath == "" {
		if inputPath == "stdin" {
			// Write to stdout
			logger.Verbose("Writing compressed data to stdout")
			_, err := stdout.Write(compressed)
			return err
		}
		// Auto-generate filename
		outputPath = inputPath + ".huff"
	}

	// Write to file
	logger.Verbose("Writing compressed file", "path", outputPath)
	if err := os.WriteFile(outputPath, compressed, 0644); err != nil {
		return fmt.Errorf("failed to write output: %w", err)
	}

	// Show stats
	originalSize := len(data)
	compressedSize := len(compressed)
	ratio := float64(compressedSize) / float64(originalSize) * 100

	fmt.Fprintf(stdout, "Compressed to: %s\n", outputPath)
	fmt.Fprintf(stdout, "Original size: %d bytes\n", originalSize)
	fmt.Fprintf(stdout, "Compressed size: %d bytes\n", compressedSize)
	fmt.Fprintf(stdout, "Compression ratio: %.2f%%\n", ratio)

	return nil
}

// doDecompress performs decompression
func (c *Command) doDecompress(data []byte, inputPath, outputPath string, stdout io.Writer) error {
	logger.Debug("Starting decompression", "input_size", len(data))

	// Decompress
	decompressed, err := Decompress(data)
	if err != nil {
		return fmt.Errorf("decompression failed: %w", err)
	}

	// Determine output
	if outputPath == "" {
		if inputPath == "stdin" {
			// Write to stdout
			logger.Verbose("Writing decompressed data to stdout")
			_, err := stdout.Write([]byte(decompressed))
			return err
		}
		// Auto-generate filename
		outputPath = strings.TrimSuffix(inputPath, ".huff")
		if outputPath == inputPath {
			outputPath = inputPath + ".decompressed"
		}
	}

	// Write to file
	logger.Verbose("Writing decompressed file", "path", outputPath)
	if err := os.WriteFile(outputPath, []byte(decompressed), 0644); err != nil {
		return fmt.Errorf("failed to write output: %w", err)
	}

	// Show stats
	fmt.Fprintf(stdout, "Decompressed to: %s\n", outputPath)
	fmt.Fprintf(stdout, "Output size: %d bytes\n", len(decompressed))

	return nil
}
