package heap

import (
    stdheap "container/heap"
    "testing"
)

// KthLargest streams values and reports the kth largest seen so far.
// Approach: min-heap of size k; after each push, drain until the heap is back down to k
// time: O(n log k), space: O(k) — the heap is capped at k elements, so its root is the kth largest
type KthLargest struct {
    minHeap *IntHeap
    k       int
}

func ConstructorKthLargest(k int, nums []int) *KthLargest {
    h := IntHeap(append([]int{}, nums...))
    stdheap.Init(&h)

    kl := &KthLargest{minHeap: &h, k: k}
    kl.maintainHeap()
    return kl
}

func (kl *KthLargest) Add(val int) int {
    stdheap.Push(kl.minHeap, val)
    kl.maintainHeap()
    return kl.minHeap.Peek()
}

func (kl *KthLargest) maintainHeap() {
    for kl.minHeap.Len() > kl.k {
        stdheap.Pop(kl.minHeap)
    }
}

// KthLargestOptimal answers the same question with one heap operation per add.
// Approach: same size-k min-heap, but once it is full replace the root in place instead of
// pushing then popping — half the heap reshuffling for the same asymptotic cost
// time: O(n log k), space: O(k)
type KthLargestOptimal struct {
    minHeap *IntHeap
    k       int
}

func ConstructorKthLargestOptimal(k int, nums []int) *KthLargestOptimal {
    h := IntHeap(append([]int{}, nums...))
    stdheap.Init(&h)

    for h.Len() > k {
        stdheap.Pop(&h)
    }

    return &KthLargestOptimal{minHeap: &h, k: k}
}

func (kl *KthLargestOptimal) Add(val int) int {
    if kl.minHeap.Len() < kl.k {
        stdheap.Push(kl.minHeap, val)
    } else if val > kl.minHeap.Peek() {
        // heapreplace: overwrite the root, then sift it down — one operation, no resize
        (*kl.minHeap)[0] = val
        stdheap.Fix(kl.minHeap, 0)
    }

    return kl.minHeap.Peek()
}

func TestKthLargestStream(t *testing.T) {
    cases := []struct {
        k     int
        nums  []int
        adds  []int
        wants []int
    }{
        {3, []int{4, 5, 8, 2}, []int{3, 5, 10, 9, 4}, []int{4, 5, 5, 8, 8}},
        {4, []int{7, 7, 7, 7, 8, 3}, []int{2, 10, 9, 9}, []int{7, 7, 7, 8}},
    }

    for _, c := range cases {
        kl := ConstructorKthLargest(c.k, c.nums)
        klOptimal := ConstructorKthLargestOptimal(c.k, c.nums)

        for i, val := range c.adds {
            if got := kl.Add(val); got != c.wants[i] {
                t.Fatalf("KthLargest.Add(%d) = %d, want %d", val, got, c.wants[i])
            }
            if got := klOptimal.Add(val); got != c.wants[i] {
                t.Fatalf("KthLargestOptimal.Add(%d) = %d, want %d", val, got, c.wants[i])
            }
        }
    }
}
