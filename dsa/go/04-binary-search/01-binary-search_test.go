package binary_search

import "testing"

// binarySearch returns the index of target in a sorted slice or -1.
func binarySearch(nums []int, target int) int {
    l, r := 0, len(nums)-1
    for l <= r {
        m := l + (r-l)/2
        if nums[m] == target {
            return m
        } else if nums[m] < target {
            l = m + 1
        } else {
            r = m - 1
        }
    }
    return -1
}

func TestBinarySearch(t *testing.T) {
    cases := []struct{ nums []int; target int; want int }{
        {[]int{-1,0,3,5,9,12}, 9, 4},
        {[]int{-1,0,3,5,9,12}, 2, -1},
        {[]int{1}, 1, 0},
        {[]int{}, 1, -1},
    }
    for _, c := range cases {
        if got := binarySearch(c.nums, c.target); got != c.want {
            t.Fatalf("binarySearch(%v, %d) = %d; want %d", c.nums, c.target, got, c.want)
        }
    }
}

