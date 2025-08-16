package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/drewslam/dscomp/pkg/huffman"
)

func runFile(c *huffman.Compressor) error {
	data, err := os.ReadFile(c.InPath)
	if err != nil {
		return fmt.Errorf("Failed to read file: %v", err)
	}
	return c.Compress(data, len(data))
}

func main() {
	// CLI Flags
	mode := flag.String("mode", "compress", "Mode: 'compress' or 'decompress'.")
	inputFile := flag.String("input", "", "Path to the input file.")
	outputFile := flag.String("output", "", "Path to the output file (optional).")

	flag.Parse()

	if *inputFile == "" {
		fmt.Fprintf(os.Stderr, "Error: An input file must be specified.\n")
		flag.Usage()
		os.Exit(1)
	}

	output := *outputFile
	if output == "" {
		if *mode == "compress" {
			output = *inputFile + ".huff"
		} else {
			output = *inputFile + ".out"
		}
	}

	comp := &huffman.Compressor{
		InPath:     *inputFile,
		OutPath:    output,
		Dictionary: make(map[byte]int),
	}

	switch *mode {
	case "compress":
		err := runFile(comp)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Compression error: %v\n", err)
		}
	case "decompress":
		err := comp.Decompress()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Decompression error: %v\n", err)
		}
	default:
		fmt.Fprintf(os.Stderr, "Invalid mode. Use 'compress' or 'decompress'.\n")
		flag.Usage()
	}

}
