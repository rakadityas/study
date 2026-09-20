package sort

import "testing"

// sortArray sorts nums in place with merge sort and returns it.
// Approach: split the range in half, sort each half recursively, then merge the two sorted halves back
// time: O(n log n), space: O(n) — log n levels of recursion, each doing O(n) merge work into temp slices
func sortArray(nums []int) []int {
    mergeSort(nums, 0, len(nums)-1)
    return nums
}

func mergeSort(nums []int, left, right int) {
    if left >= right { return }

    mid := (right + left) / 2
    mergeSort(nums, left, mid)
    mergeSort(nums, mid+1, right)

    merge(nums, left, mid, right)
}

func merge(nums []int, left, mid, right int) {
    leftNums := append([]int{}, nums[left:mid+1]...)
    rightNums := append([]int{}, nums[mid+1:right+1]...)
    leftIdx, rightIdx := 0, 0

    numsIdx := left

    for leftIdx < len(leftNums) && rightIdx < len(rightNums) {
        if leftNums[leftIdx] < rightNums[rightIdx] {
            nums[numsIdx] = leftNums[leftIdx]
            leftIdx++
        } else {
            nums[numsIdx] = rightNums[rightIdx]
            rightIdx++
        }
        numsIdx++
    }

    for leftIdx < len(leftNums) {
        nums[numsIdx] = leftNums[leftIdx]
        leftIdx++
        numsIdx++
    }

    for rightIdx < len(rightNums) {
        nums[numsIdx] = rightNums[rightIdx]
        rightIdx++
        numsIdx++
    }
}

func TestSortArray(t *testing.T) {
    cases := []struct {
        in   []int
        want []int
    }{
        {[]int{5, 2, 3, 1}, []int{1, 2, 3, 5}},
        {[]int{5, 1, 1, 2, 0, 0}, []int{0, 0, 1, 1, 2, 5}},
        {[]int{}, []int{}},
        {[]int{1}, []int{1}},
    }

    for _, c := range cases {
        got := sortArray(append([]int{}, c.in...))
        if len(got) != len(c.want) { t.Fatalf("sortArray(%v) = %v, want %v", c.in, got, c.want) }
        for i := range got {
            if got[i] != c.want[i] { t.Fatalf("sortArray(%v) = %v, want %v", c.in, got, c.want) }
        }
    }
}
