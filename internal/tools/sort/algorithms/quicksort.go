package algorithms

// QuickSort implements the Sorter interface using quicksort algorithm
type QuickSort struct{}

/*
	Optimize QuickSort
	Add these improvements:

	Median-of-three pivot (first, middle, last) : already done
	Random pivot (avoid worst case)
	Insertion sort for small arrays (< 10 elements)
	Tail recursion optimization
*/

// []string is passed by ref
// int is passed as copy
/*
get the feel of it.
we move up if we are sure the matched value is less than pivot that's it
Direction of iteration must match direction of placement.
Forward iteration → place elements forward (left)
Backward iteration → place elements backward (right)
*/
func (q QuickSort) partition(arr []string, l, r int) int {
	pivot := q.getPivot(l, r)
	pivotValue := arr[pivot]
	arr[r], arr[pivot] = arr[pivot], arr[r]

	// r has my pivot value
	pivotPosition := l
	for i := l; i < r; i++ {
		if arr[i] < pivotValue {
			arr[i], arr[pivotPosition] = arr[pivotPosition], arr[i]
			pivotPosition++
		}
	}

	// r still has my pivot value
	// Place pivot at its final sorted position
	arr[r], arr[pivotPosition] = arr[pivotPosition], arr[r]
	return pivotPosition // Pivot's final sorted position
}

// pass by ref
func (q QuickSort) quickSort(arr []string, l, r int) {
	if l >= r {
		return
	}

	pivot := q.partition(arr, l, r)
	q.quickSort(arr, l, pivot-1)
	q.quickSort(arr, pivot+1, r)
}

func (q QuickSort) getPivot(l, r int) int {
	return l + (r-l)/2
}

// NOTE:choose a pivot that divides the array in half
// medium of 3 (1st, middle and last element)
// QuickSort sorts a slice of strings in-place using the quicksort algorithm.
// Average: O(n log n), Worst: O(n²), Space: O(log n)
// Sort sorts the array in-place and returns it
func (q QuickSort) Sort(arr []string) []string {
	if len(arr) < 2 {
		return arr
	}
	q.quickSort(arr, 0, len(arr)-1)
	return arr
}
