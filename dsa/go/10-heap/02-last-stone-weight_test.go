package heap

import (
    stdheap "container/heap"
    "testing"
)

// lastStoneWeight repeatedly smashes the two heaviest stones together and returns what is left.
// Approach: max-heap of stone weights; pop two, push back the difference when they differ.
// Go lets the comparator define a max-heap directly, so no value negation is needed.
// time: O(n log n), space: O(n) — O(n) to build the heap, O(log n) per smash
func lastStoneWeight(stones []int) int {
    h := MaxHeap(append([]int{}, stones...))
    stdheap.Init(&h)

    for h.Len() > 1 {
        stoneOne := stdheap.Pop(&h).(int)
        stoneTwo := stdheap.Pop(&h).(int)

        if stoneOne > stoneTwo {
            stdheap.Push(&h, stoneOne-stoneTwo)
        }
    }

    if h.Len() == 0 { return 0 }

    return h[0]
}

func TestLastStoneWeight(t *testing.T) {
    cases := []struct {
        in   []int
        want int
    }{
        {[]int{2, 7, 4, 1, 8, 1}, 1},
        {[]int{1}, 1},
        {[]int{1, 1}, 0},
        {[]int{}, 0},
    }

    for _, c := range cases {
        if got := lastStoneWeight(c.in); got != c.want {
            t.Fatalf("lastStoneWeight(%v) = %d, want %d", c.in, got, c.want)
        }
    }
}
