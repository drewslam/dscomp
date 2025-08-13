package huffman

// Base node
type Node interface {
	IsLeaf() bool
	Weight() int
}

// Leaf node
type Leaf struct {
	element byte
	weight  int
}

func NewLeafNode(el byte, wt int) *Leaf {
	return &Leaf{
		element: el,
		weight:  wt,
	}
}

func (l *Leaf) Value() byte {
	return l.element
}

func (l *Leaf) Weight() int {
	return l.weight
}

func (l *Leaf) IsLeaf() bool {
	return true
}

// Internal node
type Internal struct {
	weight int
	right  Node
	left   Node
}

func NewInternalNode(wt int, r, l Node) *Internal {
	return &Internal{
		weight: wt,
		right:  r,
		left:   l,
	}
}

func (i *Internal) Left() Node {
	return i.left
}

func (i *Internal) Right() Node {
	return i.right
}

func (i *Internal) Weight() int {
	return i.weight
}

func (i *Internal) IsLeaf() bool {
	return false
}
