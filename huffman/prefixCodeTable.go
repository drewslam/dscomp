package huffman

type CodeTable map[byte]string

func GenerateCodeTable(root Node) CodeTable {
	codeTable := make(CodeTable)
	if root == nil {
		return codeTable
	}

	generateCodes(root, []byte{}, codeTable)
	return codeTable
}

func generateCodes(node Node, prefix []byte, codeTable CodeTable) {
	switch n := node.(type) {
	case *Leaf:
		codeTable[n.element] = string(prefix)
	case *Internal:
		generateCodes(n.Left(), append(prefix, '0'), codeTable)
		generateCodes(n.Right(), append(prefix, '1'), codeTable)
	}
}
