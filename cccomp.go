package main

import (
	"container/heap"
	"fmt"
	"os"

	"github.com/drewslam/cccomp/huffman"
)

type Compressor struct {
	dictionary map[byte]int
}

func (c *Compressor) compress(source string) error {
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
		leaf := huffman.NewLeafNode(byt, count)
		tempTree := &huffman.Tree{}
		tempTree.SetRoot(leaf)
		heap.Push(tempHeap, tempTree)
	}

	if tempHeap.Len() == 0 {
		return fmt.Errorf("No valid characters found.")
	}

	huffmanTree := huffman.BuildTree(*tempHeap)

	if huffmanTree == nil || huffmanTree.Root() == nil {
		return fmt.Errorf("Failed to build Huffman tree.")
	}

	codeTable := huffman.GenerateCodeTable(huffmanTree.Root())
	for char, code := range codeTable {
		fmt.Printf("Character '%c' (ASCII %d): %s\n", char, char, code)
	}

	return nil
}

func (c *Compressor) runFile(path string) error {
	bytes, err := os.ReadFile(path)
	if err != nil {
		return fmt.Errorf("Failed to read file: %v", err)
	}
	return c.compress(string(bytes))
}

func main() {
	compressor := &Compressor{
		dictionary: make(map[byte]int),
	}

	switch len(os.Args) {
	case 2:
		if err := compressor.runFile(os.Args[1]); err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		}
	default:
		fmt.Fprintf(os.Stderr, "Usage: %s <filename>\n", os.Args[0])
	}

}
