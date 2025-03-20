package main

import (
	"container/heap"
	"flag"
	"fmt"
	"os"

	"github.com/drewslam/cccomp/huffman"
)

type Compressor struct {
	dictionary map[byte]int
}

func (c *Compressor) compress(filename, source, outPath string, originalSize int) error {
	if c.dictionary == nil {
		c.dictionary = make(map[byte]int)
	}

	for i := 0; i < len(source); i++ {
		c.dictionary[source[i]]++
	}

	if len(c.dictionary) == 0 {
		return fmt.Errorf("Empty source file.")
	}

	tempHeap := &huffman.TreeHeap{}
	heap.Init(tempHeap)

	for byt, count := range c.dictionary {
		if count <= 0 {
			continue
		}
		tempTree := huffman.NewTree(count, byt)
		// leaf := huffman.NewLeafNode(byt, count)
		// tempTree := &huffman.Tree{}
		// tempTree.SetRoot(leaf)
		heap.Push(tempHeap, tempTree)
	}

	if tempHeap.Len() == 0 {
		return fmt.Errorf("No valid characters found.")
	}

	huffmanTree := huffman.BuildTree(tempHeap)

	if huffmanTree == nil || huffmanTree.Root() == nil {
		return fmt.Errorf("Failed to build Huffman tree.")
	}

	codeTable := huffman.GenerateCodeTable(huffmanTree.Root())
	// var encodedData []byte
	var bitBuffer byte
	var bitCount uint8
	encodedData := make([]byte, 0, (len(source) * 8 / 8))

	for i := 0; i < len(source); i++ {
		char := source[i]
		if code, exists := codeTable[char]; exists {
			for _, bit := range code {
				bitBuffer <<= 1
				if bit == '1' {
					bitBuffer |= 1
				}
				bitCount++

				if bitCount == 8 {
					encodedData = append(encodedData, bitBuffer)
					bitBuffer = 0
					bitCount = 0
				}
			}
		}
	}

	if bitCount > 0 {
		encodedData = append(encodedData, bitBuffer<<(8-bitCount))
	}
	/*	var bitString strings.Builder
		bitString.Grow(len(source) * 8)

		for i := 0; i < len(source); i++ {
			char := source[i]
			if code, exists := codeTable[char]; exists {
				bitString.WriteString(code)
			}
		}

		bitStringValue := bitString.String()

		for i := 0; i < len(bitStringValue); i += 8 {
			end := i + 8
			if end > len(bitStringValue) {
				end = len(bitStringValue)
			}
			bits := bitStringValue[i:end]
			for len(bits) < 8 {
				bits += "0"
			}

			var b byte
			for j := 0; j < 8; j++ {
				if j < len(bits) && bits[j] == '1' {
					b |= 1 << (7 - j)
				}
			}
			encodedData = append(encodedData, b)
		}
		/*for char, code := range codeTable {
			fmt.Printf("Character '%q' (ASCII %d): %s\n", char, char, code)
		}*/

	if codeTable != nil {
		for char, code := range codeTable {
			fmt.Printf("Compression: Char '%c' (%d) -> %s\n", char, char, code)
		}
		err := huffman.WriteCompressedFile(outPath, codeTable, encodedData, c.dictionary, originalSize)
		if err != nil {
			return fmt.Errorf("Failed to write compressed file: %v\n", err)
		}
	}

	fmt.Println("Successfully wrote compressed file:", outPath)
	return nil
}

func (c *Compressor) decompress(inPath, outPath string) error {
	codeTable, encodedData, originalSize, err := huffman.ReadCompressedFile(inPath)
	if err != nil {
		return fmt.Errorf("Failed to read compressed file: %v", err)
	}

	reverseCodeTable := make(map[string]byte)
	for char, code := range codeTable {
		fmt.Printf("Compression: Char '%c' (%d) -> %s\n", char, char, code)
		reverseCodeTable[code] = char
	}

	decodedData := make([]byte, 0, originalSize)
	currentCode := ""
	for _, byteVal := range encodedData {
		for i := 7; i >= 0; i-- {
			bit := (byteVal >> i) & 1
			currentCode += string('0' + bit)

			if char, exists := reverseCodeTable[currentCode]; exists {
				decodedData = append(decodedData, char)
				currentCode = ""
			}
		}
	}

	fmt.Printf("Expected output size: %d, Decoded bytes: %d\n", originalSize, len(decodedData))

	err = os.WriteFile(outPath, decodedData, 0644)
	if err != nil {
		return fmt.Errorf("Failed to write decompressed file: %v", err)
	}

	return nil
}

func (c *Compressor) runFile(inPath, outPath string) error {
	bytes, err := os.ReadFile(inPath)
	if err != nil {
		return fmt.Errorf("Failed to read file: %v", err)
	}
	return c.compress(inPath, string(bytes), outPath, len(bytes))
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

	compressor := &Compressor{
		dictionary: make(map[byte]int),
	}

	switch *mode {
	case "compress":
		err := compressor.runFile(*inputFile, output)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Compression error: %v\n", err)
		}
	case "decompress":
		err := compressor.decompress(*inputFile, output)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Decompression error: %v\n", err)
		}
	default:
		fmt.Fprintf(os.Stderr, "Invalid mode. Use 'compress' or 'decompress'.\n")
		flag.Usage()
	}

}
