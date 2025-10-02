package huffman

import (
	"cli-t/internal/command"
	"cli-t/internal/shared/file"
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
	return "[OPTIONS] [FILE]"
}

func (c *Command) Description() string {
	return "Compress and decompress files using Huffman coding"
}

func (c *Command) ValidateArgs(args []string) error {
	if len(args) > 1 {
		return fmt.Errorf("too many arguments, expected 0 or 1 file")
	}
	return nil
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
			Usage:     "Output file path",
			Type:      "string",
			Default:   "",
		},
	}
}

func (c *Command) Execute(ctx context.Context, args *command.Args) error {
	decompress, _ := args.Flags["decompress"].(bool)
	outputPath, _ := args.Flags["output"].(string)

	var inputPath string
	var content string
	var contentBytes []byte
	var err error

	// Determine input source
	if len(args.Positional) == 0 {
		// Read from stdin
		if decompress {
			// For decompress, read as bytes
			contentBytes, err = file.ReadBytes("/dev/stdin")
		} else {
			content, err = file.ReadFrom(args.Stdin)
		}
		if err != nil {
			return fmt.Errorf("error reading from stdin: %w", err)
		}
		inputPath = "stdin"
	} else {
		inputPath = args.Positional[0]
		if decompress {
			contentBytes, err = file.ReadBytes(inputPath)
		} else {
			content, err = file.ReadContent(inputPath)
		}
		if err != nil {
			return fmt.Errorf("error reading file: %w", err)
		}
	}

	if decompress {
		// Decompress
		decompressed, err := Decompress(contentBytes)
		if err != nil {
			return fmt.Errorf("decompression failed: %w", err)
		}

		// Determine output path
		if outputPath == "" {
			if inputPath == "stdin" {
				outputPath = "output.txt"
			} else {
				// Remove .huff extension or add .decompressed
				outputPath = strings.TrimSuffix(inputPath, ".huff")
				if outputPath == inputPath {
					outputPath = inputPath + ".decompressed"
				}
			}
		}

		// Write output
		if err := os.WriteFile(outputPath, []byte(decompressed), 0644); err != nil {
			return fmt.Errorf("failed to write output: %w", err)
		}

		fmt.Fprintf(args.Stdout, "Decompressed to: %s\n", outputPath)
		fmt.Fprintf(args.Stdout, "Output size: %d bytes\n", len(decompressed))

	} else {
		// Compress
		compressed, err := Compress(content)
		if err != nil {
			return fmt.Errorf("compression failed: %w", err)
		}

		// Determine output path
		if outputPath == "" {
			if inputPath == "stdin" {
				outputPath = "output.huff"
			} else {
				outputPath = inputPath + ".huff"
			}
		}

		// Write output
		if err := os.WriteFile(outputPath, compressed, 0644); err != nil {
			return fmt.Errorf("failed to write output: %w", err)
		}

		originalSize := len(content)
		compressedSize := len(compressed)
		ratio := float64(compressedSize) / float64(originalSize) * 100

		fmt.Fprintf(args.Stdout, "Compressed to: %s\n", outputPath)
		fmt.Fprintf(args.Stdout, "Original size: %d bytes\n", originalSize)
		fmt.Fprintf(args.Stdout, "Compressed size: %d bytes\n", compressedSize)
		fmt.Fprintf(args.Stdout, "Compression ratio: %.2f%%\n", ratio)
	}

	return nil
}
