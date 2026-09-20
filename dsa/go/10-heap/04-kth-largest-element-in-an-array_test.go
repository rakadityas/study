package heap

import (
    stdheap "container/heap"
    "testing"
)

// findKthLargest returns the kth largest value in nums.
// Approach: heapify the whole array as a max-heap, then pop k-1 times to reach the kth largest
// time: O(n log n), space: O(n) — heapify is O(n), but the k pops cost O(log n) each
func findKthLargest(nums []int, k int) int {
    h := MaxHeap(append([]int{}, nums...))
    stdheap.Init(&h)

    for k > 1 {
        stdheap.Pop(&h)
        k--
    }

    return stdheap.Pop(&h).(int)
}

// findKthLargestOptimal returns the same value while touching only k elements of memory.
// Approach: maintain a min-heap capped at k; once n elements have streamed through, the root
// is the kth largest
// time: O(n log k), space: O(k) — the heap never exceeds k elements
func findKthLargestOptimal(nums []int, k int) int {
    h := &IntHeap{}

    for i := range nums {
        stdheap.Push(h, nums[i])
        if h.Len() > k {
            stdheap.Pop(h)
        }
    }

    return h.Peek()
}

func TestFindKthLargest(t *testing.T) {
    cases := []struct {
        nums []int
        k    int
        want int
    }{
        {[]int{3, 2, 1, 5, 6, 4}, 2, 5},
        {[]int{3, 2, 3, 1, 2, 4, 5, 5, 6}, 4, 4},
        {[]int{1}, 1, 1},
    }

    for _, c := range cases {
        if got := findKthLargest(c.nums, c.k); got != c.want {
            t.Fatalf("findKthLargest(%v, %d) = %d, want %d", c.nums, c.k, got, c.want)
        }
        if got := findKthLargestOptimal(c.nums, c.k); got != c.want {
            t.Fatalf("findKthLargestOptimal(%v, %d) = %d, want %d", c.nums, c.k, got, c.want)
        }
    }
}
