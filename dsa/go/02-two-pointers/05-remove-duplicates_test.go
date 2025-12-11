package two_pointers

import "testing"

// removeDuplicates removes duplicates from sorted array and returns new length.
func removeDuplicates(nums []int) int {
    if len(nums) == 0 { return 0 }
    k := 1
    for i := 1; i < len(nums); i++ {
        if nums[i] != nums[k-1] { nums[k] = nums[i]; k++ }
    }
    return k
}

func TestRemoveDuplicates(t *testing.T) {
    nums := []int{1,1,2}
    k := removeDuplicates(nums)
    if k != 2 { t.Fatalf("expected 2, got %d", k) }
}

