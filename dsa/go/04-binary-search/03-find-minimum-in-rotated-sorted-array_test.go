package binary_search

import "testing"

func findMin(nums []int) int {
    l, r := 0, len(nums)-1
    for l < r {
        m := l + (r-l)/2
        if nums[m] > nums[r] { l = m + 1 } else { r = m }
    }
    return nums[l]
}

func TestFindMin(t *testing.T) {
    if findMin([]int{3,4,5,1,2}) != 1 { t.Fatalf("expected 1") }
    if findMin([]int{4,5,6,7,0,1,2}) != 0 { t.Fatalf("expected 0") }
}

