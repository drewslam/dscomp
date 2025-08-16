package huffman

import (
	"container/heap"
	"fmt"
	"os"
)

type Compressor struct {
	InPath     string
	OutPath    string
	Dictionary map[byte]int
}

func (c *Compressor) Compress(source []byte, originalSize int) error {
	if c.Dictionary == nil {
		c.Dictionary = make(map[byte]int)
	}

	for _, b := range source {
		c.Dictionary[b]++
	}

	if len(c.Dictionary) == 0 {
		return fmt.Errorf("Empty source file.")
	}

	tempHeap := &TreeHeap{}
	heap.Init(tempHeap)

	for byt, count := range c.Dictionary {
		if count <= 0 {
			continue
		}
		tempTree := NewTree(count, byt)
		heap.Push(tempHeap, tempTree)
	}

	if tempHeap.Len() == 0 {
		return fmt.Errorf("No valid characters found.")
	}

	huffmanTree := BuildTree(tempHeap)

	if huffmanTree == nil || huffmanTree.Root() == nil {
		return fmt.Errorf("Failed to build Huffman tree.")
	}

	codeTable := GenerateCodeTable(huffmanTree.Root())
	var bitBuffer byte
	var bitCount uint8
	encodedData := make([]byte, 0, (len(source) * 8 / 8))

	for _, b := range source {
		char := b
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

	if codeTable != nil {
		for char, code := range codeTable {
			fmt.Printf("Compression: Char '%c' (%d) -> %s\n", char, char, code)
		}
		err := WriteCompressedFile(c.OutPath, codeTable, encodedData, c.Dictionary, originalSize)
		if err != nil {
			return fmt.Errorf("Failed to write compressed file: %v\n", err)
		}
	}

	fmt.Println("Successfully wrote compressed file:", c.OutPath)
	return nil
}

func (c *Compressor) Decompress() error {
	codeTable, encodedData, originalSize, err := ReadCompressedFile(c.InPath)
	if err != nil {
		return fmt.Errorf("Failed to read compressed file: %v", err)
	}

	// Single unique character
	if len(codeTable) == 1 {
		var singleChar byte
		for char := range codeTable {
			singleChar = char
			break
		}
		decodedData := make([]byte, originalSize)
		for i := 0; i < originalSize; i++ {
			decodedData[i] = singleChar
		}

		err = os.WriteFile(c.OutPath, decodedData, 0644)
		if err != nil {
			return fmt.Errorf("Failed to write decompressed file: %v", err)
		}
		return nil
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
			// Stop processing if we've decoded enough characters
			if len(decodedData) >= originalSize {
				break
			}

			bit := (byteVal >> i) & 1
			currentCode += string('0' + bit)

			if char, exists := reverseCodeTable[currentCode]; exists {
				decodedData = append(decodedData, char)
				currentCode = ""
			}
		}

		if len(decodedData) >= originalSize {
			break
		}
	}

	// Verify we got the expected amount of data
	if len(decodedData) != originalSize {
		return fmt.Errorf("Decompression mismatch: expected %d bytes, got %d", originalSize, len(decodedData))
	}

	err = os.WriteFile(c.OutPath, decodedData, 0644)
	if err != nil {
		return fmt.Errorf("Failed to write decompressed file: %v", err)
	}

	fmt.Printf("Successfully decompressed %d bytes to: %s\n", len(decodedData), c.OutPath)
	return nil
}
