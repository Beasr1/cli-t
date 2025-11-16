package algorithms

type RadixSort struct{}

// Sort sorts strings using LSD radix sort
// Time: O(n*k) where k is max string length
// Space: O(n+k)
func (r RadixSort) Sort(arr []string) []string {
	if len(arr) < 2 {
		return arr
	}

	maxRounds := r.maxLen(arr)
	for i := maxRounds - 1; i >= 0; i-- {
		arr = r.countingSortByPosition(arr, i)
	}

	return arr
}

// Find maximum string length
func (r RadixSort) maxLen(arr []string) int {
	max := 0
	for _, s := range arr {
		if len(s) > max {
			max = len(s)
		}
	}
	return max
}

// Counting sort by character at position 'pos'
/*
counting sorts works by
- counting occurances
- aggregates the values

- now aggregate values helps us finding the start positions of each value (which we can track of when to insert)
- (find start positions (then traverse from i=0) or end index (traverse from back) both will work)

- this will allow us to sort in (O(n+k+n)) : k is base

ASCII is base of 256

00 1 222 3
2. 1. 3. 1.
2. 3. 6. 7.

3 ends at index (7-1)
*/
func (r RadixSort) countingSortByPosition(arr []string, pos int) []string {
	// ascii can range to 0-255 : 256 : 2^8

	n := len(arr)
	output := make([]string, n)
	count := make([]int, 256) // buckets

	// count
	for _, s := range arr {
		ch := r.charAt(s, pos)
		count[ch]++
	}

	// Calculate cumulative counts
	for i := 1; i < 256; i++ {
		count[i] += count[i-1]
	}

	for i := n - 1; i >= 0; i-- {
		ch := r.charAt(arr[i], pos)
		output[count[ch]-1] = arr[i]
		count[ch]-- // update the end index
	}

	return output
}

// Get character at position, or special "before" value
func (r RadixSort) charAt(s string, pos int) int {
	if pos >= len(s) {
		return 0 // Before 'a'
	}
	return int(s[pos]) // ASCII value
}
