package uniq

import (
	"bufio"
	"cli-t/internal/shared/io"
	"cli-t/internal/shared/logger"
	"fmt"
)

type Options struct {
	Count    bool // -c flag
	Repeated bool // -d flag
	Unique   bool // -u flag
}

// Process reads input and writes deduplicated output
func Process(reader io.Reader, writer io.Writer, opts Options) error {
	scanner := bufio.NewScanner(reader)

	previousLine := ""
	count := 0

	for scanner.Scan() {
		currentLine := scanner.Text()
		logger.Debug("line", "line", previousLine)

		if currentLine == previousLine {
			// Duplicate - just count it
			count++
		} else {
			// New line detected!
			// Should we output the PREVIOUS line?
			if previousLine != "" {
				// output previousLine
				if err := outputLine(writer, previousLine, count, opts); err != nil {
					return err
				}
			}
			// Reset state for new line
			previousLine = currentLine
			count = 1
		}
	}

	// What happens after the loop?
	if previousLine != "" {
		// output previousLine
		if err := outputLine(writer, previousLine, count, opts); err != nil {
			return err
		}
	}

	return scanner.Err()
}

func outputLine(writer io.Writer, line string, count int, opts Options) error {
	// Determine if we should print based on flags
	// Format the output based on opts.Count
	// Write to writer
	shouldPrint := (opts.Repeated && count > 1) || (opts.Unique && count == 1) || (!opts.Repeated && !opts.Unique)
	logger.Debug("shouldprint", "shouldprint", shouldPrint)
	if !shouldPrint {
		return nil
	}

	// Step 2: Format the output
	var output string
	if opts.Count {
		output = fmt.Sprintf("%d %s", count, line)
	} else {
		output = line
	}

	// Step 3: Write it
	_, err := fmt.Fprintf(writer, "%s\n", output)
	return err

}
