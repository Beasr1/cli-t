package algorithms

type MergeSort struct{}

// MergeSort sorts a slice of strings using merge sort algorithm.
// Time: O(n log n), Space: O(n), Stable: Yes
func (m MergeSort) Sort(arr []string) []string {
	return m.mergeSort(arr)
}

func (m MergeSort) mergeSort(arr []string) []string {
	if len(arr) < 2 {
		return arr
	}

	n := len(arr)
	pivot := m.getPivot(0, n-1) // pivot is the index we should go to
	// pivot should not be less than 0 since infinite recur
	left := m.mergeSort(arr[:pivot+1])  // 0 - pivot
	right := m.mergeSort(arr[pivot+1:]) // pivot+1 - (n-1)
	return m.merge(left, right)
}

// left and right are sorted arrays
func (m MergeSort) merge(left, right []string) []string {
	merged := make([]string, 0, len(left)+len(right)) // init for fast af allocation

	// stable sort : equal elements maintain order
	l, r := 0, 0
	for l < len(left) && r < len(right) {
		if left[l] < right[r] {
			merged = append(merged, left[l])
			l++
		} else {
			merged = append(merged, right[r])
			r++
		}
	}

	// Append remaining elements
	merged = append(merged, left[l:]...)
	merged = append(merged, right[r:]...)

	return merged
}

func (m MergeSort) getPivot(l, r int) int {
	return l + (r-l)/2
}
