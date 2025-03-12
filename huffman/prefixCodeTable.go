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

func generateCodes(node Node, currentCode string, codeTable CodeTable) {
	if node.IsLeaf() {
		leaf := node.(*Leaf)
		codeTable[leaf.Value()] = currentCode
		return
	}

	internal := node.(*Internal)

	generateCodes(internal.Left(), currentCode+"0", codeTable)
	generateCodes(internal.Right(), currentCode+"1", codeTable)
}
