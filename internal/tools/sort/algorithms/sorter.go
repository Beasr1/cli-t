package algorithms

// Sorter defines the interface for all sorting algorithms.
// All sorting algorithms must implement this interface.
type Sorter interface {
	// Sort takes a slice of strings and returns a sorted slice.
	// Some implementations sort in-place, others return a new slice.
	Sort(arr []string) []string
}
