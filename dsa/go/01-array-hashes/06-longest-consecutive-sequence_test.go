package array_hashes

import "testing"

// longestConsecutive returns length of the longest consecutive sequence.
func longestConsecutive(nums []int) int {
    set := make(map[int]bool, len(nums))
    for _, n := range nums { set[n] = true }
    best := 0
    for n := range set {
        if !set[n-1] { // start of a sequence
            length := 1
            cur := n
            for set[cur+1] { cur++; length++ }
            if length > best { best = length }
        }
    }
    return best
}

func TestLongestConsecutive(t *testing.T) {
    cases := []struct{ nums []int; want int }{
        {[]int{100,4,200,1,3,2}, 4},
        {[]int{0,3,7,2,5,8,4,6,0,1}, 9},
        {[]int{}, 0},
    }
    for _, c := range cases {
        if got := longestConsecutive(c.nums); got != c.want {
            t.Fatalf("longestConsecutive(%v) = %d; want %d", c.nums, got, c.want)
        }
    }
}

