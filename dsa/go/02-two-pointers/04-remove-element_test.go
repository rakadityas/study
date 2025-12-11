package two_pointers

import "testing"

// removeElement removes all instances of val in-place and returns the new length.
func removeElement(nums []int, val int) int {
    k := 0
    for _, x := range nums {
        if x != val { nums[k] = x; k++ }
    }
    return k
}

func TestRemoveElement(t *testing.T) {
    nums := []int{3,2,2,3}
    k := removeElement(nums, 3)
    if k != 2 { t.Fatalf("expected 2, got %d", k) }
}

