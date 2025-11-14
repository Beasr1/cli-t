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
	// TODO: Implement
	return arr
}

func (m MergeSort) merge(left, right []string) []string {
	// TODO: Implement
	return nil
}

func (m MergeSort) getPivot(l, r int) int {
	return l + (r-l)/2
}
