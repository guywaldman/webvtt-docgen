package main

import (
	"flag"
	"fmt"
	"io"
	"os"
	"strings"
	"time"

	vtt_parser "github.com/guywaldman/webvtt-docgen/internal"
)

func main() {
	flag.Parse()
	srcFilePath := flag.Arg(0)
	if strings.TrimSpace(srcFilePath) == "" {
		fmt.Fprintln(os.Stderr, "ERROR: `src` argument is missing")
		printUsage(os.Stderr)
		os.Exit(1)
	}
	dstFilePath := flag.Arg(1)
	if strings.TrimSpace(dstFilePath) == "" {
		fmt.Fprintln(os.Stderr, "ERROR: `dst` argument is missing")
		printUsage(os.Stderr)
		os.Exit(1)
	}

	file, err := os.OpenFile(srcFilePath, os.O_RDONLY, 0)
	if err != nil {
		fmt.Fprintf(os.Stderr, "ERROR: error opening file at path '%s': %v\n", srcFilePath, err)
		printUsage(os.Stderr)
		os.Exit(1)
	}
	defer file.Close()

	parsedVttFile, err := vtt_parser.ParseVTT(file)
	if err != nil {
		fmt.Fprintf(os.Stderr, "ERROR: error reading from file at path '%s': %v\n", srcFilePath, err)
		printUsage(os.Stderr)
		os.Exit(1)
	}

	markdown, err := formatMarkdown(parsedVttFile)
	if err != nil {
		fmt.Fprintf(os.Stderr, "ERROR: error formatting markdown: %v\n", err)
		os.Exit(1)
	}

	os.WriteFile(dstFilePath, []byte(markdown), 0644)
	if err != nil {
		fmt.Fprintf(os.Stderr, "ERROR: error writing to file at path '%s': %v\n", dstFilePath, err)
		os.Exit(1)
	}
}

func formatMarkdown(parsedVttFile *vtt_parser.VTTFile) (string, error) {
	var sb strings.Builder

	sb.WriteString("> Generated by https://github.com/guywaldman/webvtt-docgen\n\n")

	for _, cue := range parsedVttFile.Cues {
		parsedStartTime, timeParsingErr := time.Parse("15:04:05.000", cue.Start)
		if timeParsingErr != nil {
			return "", timeParsingErr
		}
		startTime := parsedStartTime.Format("15:04:05")
		sb.WriteString(fmt.Sprintf("- **%s** (%s)  \n", cue.Speaker, startTime))
		for _, line := range strings.Split(cue.Text, "\n") {
			sb.WriteString(fmt.Sprintf("  %s\n", strings.TrimSpace(line)))
		}
	}
	return sb.String(), nil
}

func printUsage(w io.Writer) {
	fmt.Fprintln(w, "Usage: <program> <src> <dest>")
	fmt.Fprintln(w, "Example usage: <program> /path/to/transcript.vtt /path/to/output.md`)")
}
