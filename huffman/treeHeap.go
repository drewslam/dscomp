package huffman

import "container/heap"

// PriorityQueue implementation for Huffman tree
type TreeHeap []*Tree

func (t TreeHeap) Len() int           { return len(t) }
func (t TreeHeap) Less(i, j int) bool { return t[i].Weight() < t[j].Weight() }
func (t TreeHeap) Swap(i, j int)      { t[i], t[j] = t[j], t[i] }

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

func BuildTree(trees TreeHeap) *Tree {
	// Convert to heap
	heap.Init(&trees)

	// Continue until only one tree remaining
	for trees.Len() > 1 {
		tmp1 := heap.Pop(&trees).(*Tree)
		tmp2 := heap.Pop(&trees).(*Tree)

		combinedWeight := tmp1.Weight() + tmp2.Weight()
		newTree := NewTree(combinedWeight, tmp1.Root(), tmp2.Root())

		heap.Push(&trees, newTree)
	}

	if trees.Len() > 0 {
		return trees[0]
	}
	return nil
}
