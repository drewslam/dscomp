package huffman

type CodeTable map[byte]string

func GenerateCodeTable(root Node) CodeTable {
	codeTable := make(CodeTable)
	if root == nil {
		return codeTable
	}

	generateCodes(root, "", codeTable)
	return codeTable
}

func generateCodes(node Node, prefix string, codeTable CodeTable) {
	switch n := node.(type) {
	case *Leaf:
		codeTable[n.element] = prefix
		return
	case *Internal:
		generateCodes(n.Left(), prefix+"0", codeTable)
		generateCodes(n.Right(), prefix+"1", codeTable)
	}
}
