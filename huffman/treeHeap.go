package huffman

import "container/heap"

// PriorityQueue implementation for Huffman tree
type TreeHeap []*Tree

func (t TreeHeap) Len() int { return len(t) }
func (t TreeHeap) Less(i, j int) bool {
	if t[i].Weight() != t[j].Weight() {
		return t[i].Weight() < t[j].Weight()
	}
	iIsLeaf := t[i].root.IsLeaf()
	jIsLeaf := t[j].root.IsLeaf()

	if iIsLeaf && jIsLeaf {
		iLeaf := t[i].root.(*Leaf)
		jLeaf := t[j].root.(*Leaf)
		return iLeaf.Value() < jLeaf.Value()
	} else if iIsLeaf {
		return true
	} else if jIsLeaf {
		return false
	}

	return getMinLeafVal(t[i].root) < getMinLeafVal(t[j].root)
}
func (t TreeHeap) Swap(i, j int) { t[i], t[j] = t[j], t[i] }

func (t *TreeHeap) Push(x interface{}) {
	*t = append(*t, x.(*Tree))
}

func (t *TreeHeap) Pop() interface{} {
	old := *t
	n := len(old)
	x := old[n-1]
	*t = old[0 : n-1]
	return x
}

func BuildTree(trees *TreeHeap) *Tree {
	// Convert to heap
	heap.Init(trees)

	// Continue until only one tree remaining
	for trees.Len() > 1 {
		tmp1 := heap.Pop(trees).(*Tree)
		tmp2 := heap.Pop(trees).(*Tree)

		var left, right Node

		if tmp1.Weight() < tmp2.Weight() {
			left = tmp1.Root()
			right = tmp2.Root()
		} else if tmp2.Weight() < tmp1.Weight() {
			left = tmp2.Root()
			right = tmp1.Root()
		} else {
			tmp1IsLeaf := tmp1.Root().IsLeaf()
			tmp2IsLeaf := tmp2.Root().IsLeaf()

			if tmp1IsLeaf && tmp2IsLeaf {
				tmp1Leaf := tmp1.Root().(*Leaf)
				tmp2Leaf := tmp2.Root().(*Leaf)

				if tmp1Leaf.Value() < tmp2Leaf.Value() {
					left = tmp1.Root()
					right = tmp2.Root()
				} else {
					left = tmp2.Root()
					right = tmp1.Root()
				}
			} else if tmp1IsLeaf {
				left = tmp1.Root()
				right = tmp2.Root()
			} else if tmp2IsLeaf {
				left = tmp2.Root()
				right = tmp1.Root()
			} else {
				if getMinLeafVal(tmp1.Root()) < getMinLeafVal(tmp2.Root()) {
					left = tmp1.Root()
					right = tmp2.Root()
				} else {
					left = tmp2.Root()
					right = tmp1.Root()
				}
			}
		}

		combinedWeight := tmp1.Weight() + tmp2.Weight()
		newTree := NewTree(combinedWeight, left, right)

		heap.Push(trees, newTree)
	}

	if trees.Len() == 0 {
		return nil
	}
	return heap.Pop(trees).(*Tree)
}

func getMinLeafVal(node Node) byte {
	switch n := node.(type) {
	case *Leaf:
		return n.Value()
	case *Internal:
		leftMin := getMinLeafVal(n.Left())
		rightMin := getMinLeafVal(n.Right())
		if leftMin < rightMin {
			return leftMin
		}
		return rightMin
	}
	return 0
}
