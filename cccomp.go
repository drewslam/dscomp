package main

import (
	"fmt"
	"os"
)

type Compressor struct {
	dictionary map[byte]int
}

func (c *Compressor) compress(source string) error {
	for i := 0; i < len(source); i++ {
		c.dictionary[source[i]]++
	}

	for byt, count := range c.dictionary {
		switch byt {
		case '\n':
			fmt.Printf("'\\n' (ASCII %d): %d\n", byt, count)
		case '\r':
			fmt.Printf("'\\r' (ASCII %d): %d\n", byt, count)
		case '\t':
			fmt.Printf("'\\t' (ASCII %d): %d\n", byt, count)
		default:
			if byt >= 32 && byt < 127 {
				fmt.Printf("%c (ASCII %d): %d\n", byt, byt, count)
			} else {
				fmt.Printf("(ASCII %d): %d\n", byt, count)
			}
		}
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
