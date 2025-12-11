package binary_search

import "testing"

func searchRotated(nums []int, target int) int {
    l, r := 0, len(nums)-1
    for l <= r {
        m := l + (r-l)/2
        if nums[m] == target { return m }
        if nums[l] <= nums[m] { // left sorted
            if nums[l] <= target && target < nums[m] { r = m - 1 } else { l = m + 1 }
        } else { // right sorted
            if nums[m] < target && target <= nums[r] { l = m + 1 } else { r = m - 1 }
        }
    }
    return -1
}

func TestSearchRotated(t *testing.T) {
    if searchRotated([]int{4,5,6,7,0,1,2}, 0) != 4 { t.Fatalf("expected 4") }
    if searchRotated([]int{4,5,6,7,0,1,2}, 3) != -1 { t.Fatalf("expected -1") }
}

