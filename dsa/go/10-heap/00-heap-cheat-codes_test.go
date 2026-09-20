package heap

import (
    stdheap "container/heap"
    "sort"
    "testing"
)

// Go has no built-in heap type — container/heap is an *algorithm* that operates on any type
// implementing heap.Interface (sort.Interface plus Push/Pop). That is the main difference from
// Python's heapq, which works directly on a list.
//
// Time complexities:
//   stdheap.Init(h)      - O(n)          (Python: heapq.heapify)
//   stdheap.Push(h, x)   - O(log n)      (Python: heapq.heappush)
//   stdheap.Pop(h)       - O(log n)      (Python: heapq.heappop)
//   stdheap.Fix(h, i)    - O(log n)      (no direct Python equivalent)
//   h[0] (peek)          - O(1)          (Python: heap[0])
//
// Space: O(n) for storage, O(1) per operation.
//
// Heap layout (same as Python): for the node at index i,
//   parent      = (i-1)/2
//   left child  = 2*i+1
//   right child = 2*i+2

// IntHeap is the min-heap used across this package. It mirrors Python's default heapq behaviour.
type IntHeap []int

func (h IntHeap) Len() int            { return len(h) }
func (h IntHeap) Less(i, j int) bool  { return h[i] < h[j] }
func (h IntHeap) Swap(i, j int)       { h[i], h[j] = h[j], h[i] }
func (h *IntHeap) Push(x any)         { *h = append(*h, x.(int)) }
func (h *IntHeap) Pop() any           { old := *h; n := len(old); v := old[n-1]; *h = old[:n-1]; return v }
func (h IntHeap) Peek() int           { return h[0] }

// MaxHeap flips the comparison instead of negating values the way the Python bank does —
// in Go the comparator is ours to define, so there is no need for the negation trick.
type MaxHeap []int

func (h MaxHeap) Len() int           { return len(h) }
func (h MaxHeap) Less(i, j int) bool { return h[i] > h[j] }
func (h MaxHeap) Swap(i, j int)      { h[i], h[j] = h[j], h[i] }
func (h *MaxHeap) Push(x any)        { *h = append(*h, x.(int)) }
func (h *MaxHeap) Pop() any          { old := *h; n := len(old); v := old[n-1]; *h = old[:n-1]; return v }

// Item is a priority-queue entry. index breaks ties in insertion order, matching the
// (priority, index, item) tuple trick used with heapq.
type Item struct {
    Value    string
    Priority int
    Index    int
}

// PriorityQueue pops the lowest Priority first, oldest insertion winning ties.
type PriorityQueue []*Item

func (pq PriorityQueue) Len() int { return len(pq) }
func (pq PriorityQueue) Less(i, j int) bool {
    if pq[i].Priority != pq[j].Priority { return pq[i].Priority < pq[j].Priority }
    return pq[i].Index < pq[j].Index
}
func (pq PriorityQueue) Swap(i, j int) { pq[i], pq[j] = pq[j], pq[i] }
func (pq *PriorityQueue) Push(x any)   { *pq = append(*pq, x.(*Item)) }
func (pq *PriorityQueue) Pop() any     { old := *pq; n := len(old); v := old[n-1]; *pq = old[:n-1]; return v }

// kLargestStreaming keeps a min-heap capped at k, so the root is always the kth largest.
// time: O(n log k), space: O(k)
func kLargestStreaming(arr []int, k int) []int {
    h := &IntHeap{}
    stdheap.Init(h)
    for _, num := range arr {
        if h.Len() < k {
            stdheap.Push(h, num)
        } else if num > h.Peek() {
            stdheap.Pop(h)
            stdheap.Push(h, num)
        }
    }

    res := append([]int{}, *h...)
    sort.Sort(sort.Reverse(sort.IntSlice(res)))
    return res
}

// cursor is one position inside one of the input lists.
type cursor struct {
    val     int
    listIdx int
    elemIdx int
}

// cursorHeap orders cursors by value, breaking ties on list index for determinism.
type cursorHeap []cursor

func (h cursorHeap) Len() int { return len(h) }
func (h cursorHeap) Less(i, j int) bool {
    if h[i].val != h[j].val { return h[i].val < h[j].val }
    return h[i].listIdx < h[j].listIdx
}
func (h cursorHeap) Swap(i, j int) { h[i], h[j] = h[j], h[i] }
func (h *cursorHeap) Push(x any)   { *h = append(*h, x.(cursor)) }
func (h *cursorHeap) Pop() any     { old := *h; n := len(old); v := old[n-1]; *h = old[:n-1]; return v }

// mergeKSortedLists merges k sorted slices by keeping one cursor per list on a heap.
// time: O(n log k), space: O(k) — n total elements across k lists
func mergeKSortedLists(lists [][]int) []int {
    h := &cursorHeap{}
    for i, lst := range lists {
        if len(lst) > 0 {
            stdheap.Push(h, cursor{lst[0], i, 0})
        }
    }

    result := []int{}
    for h.Len() > 0 {
        c := stdheap.Pop(h).(cursor)
        result = append(result, c.val)

        if c.elemIdx+1 < len(lists[c.listIdx]) {
            stdheap.Push(h, cursor{lists[c.listIdx][c.elemIdx+1], c.listIdx, c.elemIdx + 1})
        }
    }

    return result
}

func TestHeapCheatCodes(t *testing.T) {
    // min-heap: Init then pop ascending
    h := &IntHeap{3, 1, 4, 1, 5, 9, 2, 6}
    stdheap.Init(h)
    if h.Peek() != 1 { t.Fatalf("min-heap root = %d, want 1", h.Peek()) }

    stdheap.Push(h, 0)
    if got := stdheap.Pop(h).(int); got != 0 { t.Fatalf("pop = %d, want 0", got) }

    // max-heap: comparator flipped, no value negation needed
    mh := &MaxHeap{}
    stdheap.Push(mh, 3)
    stdheap.Push(mh, 1)
    stdheap.Push(mh, 4)
    if got := stdheap.Pop(mh).(int); got != 4 { t.Fatalf("max-heap pop = %d, want 4", got) }
    if (*mh)[0] != 3 { t.Fatalf("max-heap peek = %d, want 3", (*mh)[0]) }

    // priority queue: lowest priority number first
    pq := &PriorityQueue{}
    stdheap.Push(pq, &Item{Value: "task1", Priority: 3, Index: 0})
    stdheap.Push(pq, &Item{Value: "task2", Priority: 1, Index: 1})
    stdheap.Push(pq, &Item{Value: "task3", Priority: 2, Index: 2})
    if got := stdheap.Pop(pq).(*Item).Value; got != "task2" { t.Fatalf("pq pop = %q, want task2", got) }
    if got := stdheap.Pop(pq).(*Item).Value; got != "task3" { t.Fatalf("pq pop = %q, want task3", got) }

    // top-k over a stream
    want := []int{9, 6, 5}
    got := kLargestStreaming([]int{3, 1, 4, 1, 5, 9, 2, 6}, 3)
    for i := range want {
        if got[i] != want[i] { t.Fatalf("kLargestStreaming = %v, want %v", got, want) }
    }

    // merge k sorted lists
    wantMerged := []int{1, 1, 2, 3, 4, 4, 5, 6}
    gotMerged := mergeKSortedLists([][]int{{1, 4, 5}, {1, 3, 4}, {2, 6}})
    for i := range wantMerged {
        if gotMerged[i] != wantMerged[i] { t.Fatalf("mergeKSortedLists = %v, want %v", gotMerged, wantMerged) }
    }
}
