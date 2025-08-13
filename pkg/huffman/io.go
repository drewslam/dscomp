package huffman

import (
	"container/heap"
	"encoding/binary"
	"fmt"
	"io"
	"os"
)

const MagicNumber = "HFZ1"

func WriteCompressedFile(filename string, codeTable CodeTable, encodedData []byte, freqMap map[byte]int, originalSize int) (retErr error) {
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer func() {
		closeErr := file.Close()
		if retErr == nil && closeErr != nil {
			retErr = closeErr
		}

	}()

	_, err = file.Write([]byte(MagicNumber))
	if err != nil {
		return err
	}

	err = binary.Write(file, binary.BigEndian, int64(originalSize))
	if err != nil {
		return err
	}

	tableSize := int32(len(freqMap))
	err = binary.Write(file, binary.BigEndian, tableSize)
	if err != nil {
		return err
	}

	for char, freq := range freqMap {
		err = binary.Write(file, binary.BigEndian, char)
		if err != nil {
			return err
		}
		err = binary.Write(file, binary.BigEndian, int32(freq))
		if err != nil {
			return err
		}
	}

	_, err = file.Write(encodedData)
	if err != nil {
		return err
	}

	return nil
}

func ReadCompressedFile(filename string) (map[byte]string, []byte, int, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, nil, 0, err
	}
	defer file.Close()

	magic := make([]byte, len(MagicNumber))
	_, err = io.ReadFull(file, magic)
	if err != nil || string(magic) != MagicNumber {
		return nil, nil, 0, fmt.Errorf("Invalid file format.")
	}

	var originalSize int64
	err = binary.Read(file, binary.BigEndian, &originalSize)
	if err != nil {
		return nil, nil, 0, err
	}

	var tableSize int32
	err = binary.Read(file, binary.BigEndian, &tableSize)
	if err != nil {
		return nil, nil, 0, err
	}

	freqMap := make(map[byte]int)
	for i := 0; i < int(tableSize); i++ {
		var char byte
		var freq int32
		err = binary.Read(file, binary.BigEndian, &char)
		if err != nil {
			return nil, nil, 0, err
		}
		err = binary.Read(file, binary.BigEndian, &freq)
		if err != nil {
			return nil, nil, 0, err
		}
		freqMap[char] = int(freq)
	}

	tempHeap := &TreeHeap{}
	heap.Init(tempHeap)

	for char, freq := range freqMap {
		heap.Push(tempHeap, NewTree(freq, char))
	}

	huffmanTree := BuildTree(tempHeap)
	if huffmanTree == nil || huffmanTree.Root() == nil {
		return nil, nil, 0, fmt.Errorf("Failed to reconstruct Huffman tree.")
	}

	codeTable := GenerateCodeTable(huffmanTree.Root())

	encodedData, err := io.ReadAll(file)
	if err != nil {
		return nil, nil, 0, err
	}

	return codeTable, encodedData, int(originalSize), nil
}
