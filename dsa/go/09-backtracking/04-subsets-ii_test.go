package backtracking

import (
    "sort"
    "strconv"
    "strings"
    "testing"
)

// subsetsWithDupSkip returns every distinct subset of a slice that may contain duplicates.
// Approach: sort first, then on the "skip" branch advance past every copy of the current value
// so a value is never dropped twice at the same position — no dedup structure needed
// time: O(2^n), space: O(n) — at most 2^n subsets after dedup; recursion stack only
func subsetsWithDupSkip(nums []int) [][]int {
    res := [][]int{}
    nums = append([]int{}, nums...)
    sort.Ints(nums)

    var backtracking func(i int, currentSet []int)
    backtracking = func(i int, currentSet []int) {
        if i >= len(nums) {
            res = append(res, append([]int{}, currentSet...))
            return
        }

        currentSet = append(currentSet, nums[i])
        backtracking(i+1, currentSet)

        currentSet = currentSet[:len(currentSet)-1]
        for i+1 < len(nums) && nums[i] == nums[i+1] {
            i++
        }
        backtracking(i+1, currentSet)
    }

    backtracking(0, []int{})
    return res
}

// subsetsWithDupSeen returns the same subsets by generating all of them and filtering.
// Approach: generate every subset, deduplicate via a seen-set keyed on the path — easier to reason
// about but trades memory for that simplicity
// time: O(2^n), space: O(2^n) — the seen set can hold up to 2^n keys
func subsetsWithDupSeen(nums []int) [][]int {
    res := [][]int{}
    nums = append([]int{}, nums...)
    sort.Ints(nums)
    seen := map[string]bool{}

    key := func(path []int) string {
        parts := make([]string, 0, len(path))
        for _, v := range path {
            parts = append(parts, strconv.Itoa(v))
        }
        return strings.Join(parts, ",")
    }

    var backtrack func(start int, path []int)
    backtrack = func(start int, path []int) {
        k := key(path)
        if !seen[k] {
            seen[k] = true
            res = append(res, append([]int{}, path...))
        }

        for i := start; i < len(nums); i++ {
            path = append(path, nums[i])
            backtrack(i+1, path)
            path = path[:len(path)-1]
        }
    }

    backtrack(0, []int{})
    return res
}

func TestSubsetsWithDup(t *testing.T) {
    cases := []struct {
        in   []int
        want [][]int
    }{
        {[]int{1, 2, 2}, [][]int{{}, {1}, {1, 2}, {1, 2, 2}, {2}, {2, 2}}},
        {[]int{0}, [][]int{{}, {0}}},
        {[]int{4, 4, 1, 4}, [][]int{{}, {1}, {1, 4}, {1, 4, 4}, {1, 4, 4, 4}, {4}, {4, 4}, {4, 4, 4}}},
    }

    for _, c := range cases {
        if got := subsetsWithDupSkip(c.in); !equalIgnoringOrder(got, c.want) {
            t.Fatalf("subsetsWithDupSkip(%v) = %v, want %v", c.in, got, c.want)
        }
        if got := subsetsWithDupSeen(c.in); !equalIgnoringOrder(got, c.want) {
            t.Fatalf("subsetsWithDupSeen(%v) = %v, want %v", c.in, got, c.want)
        }
    }
}
