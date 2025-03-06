package huffman

// Tree implementation
type Tree struct {
	root Node
}

// Tree constructor
func NewTree(wt int, args ...interface{}) *Tree {
	switch len(args) {
	case 1:
		// Single argument case
		if val, ok := args[0].(byte); ok {
			return &Tree{
				root: NewLeafNode(val, wt),
			}
		}
	case 2:
		// Two argument case
		if right, ok := args[0].(Node); ok {
			if left, ok := args[1].(Node); ok {
				return &Tree{
					root: NewInternalNode(wt, right, left),
				}
			}
		}
	}

	return &Tree{
		root: NewLeafNode(0, 0), // Default for invalid argument format
	}
}

func (t *Tree) SetRoot(r Node) {
	t.root = r
}

func (t *Tree) Root() Node {
	return t.root
}

func (t *Tree) Weight() int {
	if t == nil || t.root == nil {
		return 0
	}
	return t.root.Weight()
}

func (t *Tree) CompareTo(u interface{}) int {
	that, ok := u.(*Tree)
	if !ok {
		return 1
	}

	if t.root.Weight() < that.Weight() {
		return -1
	} else if t.root.Weight() == that.Weight() {
		return 0
	}

	return 1
}
