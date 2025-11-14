package algorithms

type HeapSort struct{}

// I can do it inplace. create new array and put it there etc etc
// Sort sorts the array in-place using heap sort algorithm.
// Time: O(n log n), Space: O(1), Stable: No
func (h HeapSort) Sort(arr []string) []string {
	if len(arr) < 2 {
		return arr
	}

	n := len(arr)

	// Phase 1: Build max heap
	h.heapify(arr, n)

	// Phase 2: Extract elements one by one
	for i := n - 1; i > 0; i-- { // Stop at 1, not 0
		arr[0], arr[i] = arr[i], arr[0]
		h.down(arr, 0, i)
	}

	// btw I can use same array as well
	// sunce arr is currently in decending order
	// h.reverse(arr)
	return arr
}

/*
So following this code we would notice
if using minheap we create decending order array
if using max heap we create ascending order array
just switch cmp()
*/
func (h HeapSort) reverse(arr []string) {
	left, right := 0, len(arr)-1
	for {
		if left >= right {
			break
		}

		arr[left], arr[right] = arr[right], arr[left]
		left++
		right--
	}
}

// heapify restores heap property for subtree rooted at index i
// n is the size of the heap (not necessarily len(arr))
// max nlogn
func (h HeapSort) heapify(arr []string, n int) {
	for i := (n/2 - 1); i >= 0; i-- {
		// swaps if parent is greater means : we are creating max heap
		// we should keep it downing and swapping
		h.down(arr, i, n)
	}
}

// Helper functions for array indices
func (h HeapSort) leftChild(i int) int {
	return 2*i + 1
}

func (h HeapSort) rightChild(i int) int {
	return 2*i + 2
}

// swap if true
func (h HeapSort) cmp(arr []string, i, j int) bool {
	return arr[i] > arr[j]
}

// parent will have max : greater value
func (h HeapSort) down(arr []string, i0, n int) {
	i := i0
	for {
		j1 := h.leftChild(i)
		if j1 >= n || j1 < 0 { // j1 < 0 after int overflow
			break
		}

		j := j1 // left child
		if j2 := h.rightChild(i); j2 < n && h.cmp(arr, j2, j1) {
			j = j2 // = 2*i + 2  // right child
		}

		// j has the  maximum element out of the child
		if !h.cmp(arr, j, i) {
			break
		}

		// swap is j is greater than i
		arr[i], arr[j] = arr[j], arr[i]
		i = j
	}
}
