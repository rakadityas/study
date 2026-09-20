package heap

import (
    stdheap "container/heap"
    "sort"
    "testing"
)

// pointHeap is a max-heap on squared distance, so its root is the worst of the k best so far.
type pointHeap []struct {
    dist  int
    point []int
}

func (h pointHeap) Len() int           { return len(h) }
func (h pointHeap) Less(i, j int) bool { return h[i].dist > h[j].dist }
func (h pointHeap) Swap(i, j int)      { h[i], h[j] = h[j], h[i] }
func (h *pointHeap) Push(x any) {
    *h = append(*h, x.(struct {
        dist  int
        point []int
    }))
}
func (h *pointHeap) Pop() any { old := *h; n := len(old); v := old[n-1]; *h = old[:n-1]; return v }

// kClosest returns the k points nearest the origin, in no particular order.
// Approach: push every point onto a max-heap keyed on squared distance and evict the root
// whenever the heap exceeds k, so only the k closest survive. Squared distance avoids the sqrt.
// time: O(n log k), space: O(k) — the heap never holds more than k+1 points
func kClosest(points [][]int, k int) [][]int {
    h := &pointHeap{}

    for i := range points {
        dist := points[i][0]*points[i][0] + points[i][1]*points[i][1]
        stdheap.Push(h, struct {
            dist  int
            point []int
        }{dist, points[i]})

        if h.Len() > k {
            stdheap.Pop(h)
        }
    }

    result := [][]int{}
    for _, entry := range *h {
        result = append(result, entry.point)
    }

    return result
}

func TestKClosest(t *testing.T) {
    cases := []struct {
        points [][]int
        k      int
        want   [][]int
    }{
        {[][]int{{1, 3}, {-2, 2}}, 1, [][]int{{-2, 2}}},
        {[][]int{{3, 3}, {5, -1}, {-2, 4}}, 2, [][]int{{-2, 4}, {3, 3}}},
    }

    for _, c := range cases {
        got := kClosest(c.points, c.k)
        // the heap returns the k closest in arbitrary order, so sort before comparing
        sort.Slice(got, func(i, j int) bool { return got[i][0] < got[j][0] })

        if len(got) != len(c.want) {
            t.Fatalf("kClosest(%v, %d) = %v, want %v", c.points, c.k, got, c.want)
        }
        for i := range got {
            if got[i][0] != c.want[i][0] || got[i][1] != c.want[i][1] {
                t.Fatalf("kClosest(%v, %d) = %v, want %v", c.points, c.k, got, c.want)
            }
        }
    }
}
