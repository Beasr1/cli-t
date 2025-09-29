package huffman

// An HuffmanMinHeap is a min-heap of ints.
type HuffmanMinHeap []*HuffmanNode

func (h HuffmanMinHeap) Len() int           { return len(h) }
func (h HuffmanMinHeap) Less(i, j int) bool { return h[i].frequency < h[j].frequency }
func (h HuffmanMinHeap) Swap(i, j int)      { h[i], h[j] = h[j], h[i] }

func (h *HuffmanMinHeap) Push(x any) {
	// Push and Pop use pointer receivers because they modify the slice's length,
	// not just its contents.
	*h = append(*h, x.(*HuffmanNode))
}

func (h *HuffmanMinHeap) Pop() any {
	old := *h
	n := len(old)
	x := old[n-1]
	old[n-1] = nil // optimization for early garbage collection
	*h = old[0 : n-1]
	return x
}
