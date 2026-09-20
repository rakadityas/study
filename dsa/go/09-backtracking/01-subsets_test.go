package backtracking

import "testing"

// subsets returns every subset of nums.
// Approach: for each index branch twice — take nums[i], then backtrack and skip it
// time: O(n*2^n), space: O(n*2^n) — each element is independently in or out, and every subset is copied out
func subsets(nums []int) [][]int {
    res := [][]int{}

    var backtracking func(i int, currentSet []int)
    backtracking = func(i int, currentSet []int) {
        if i >= len(nums) {
            res = append(res, append([]int{}, currentSet...))
            return
        }

        currentSet = append(currentSet, nums[i])
        backtracking(i+1, currentSet)

        currentSet = currentSet[:len(currentSet)-1]
        backtracking(i+1, currentSet)
    }

    backtracking(0, []int{})
    return res
}

func TestSubsets(t *testing.T) {
    want := [][]int{{}, {1}, {2}, {1, 2}, {3}, {1, 3}, {2, 3}, {1, 2, 3}}
    if got := subsets([]int{1, 2, 3}); !equalIgnoringOrder(got, want) {
        t.Fatalf("subsets([1 2 3]) = %v, want %v", got, want)
    }

    want = [][]int{{}, {0}}
    if got := subsets([]int{0}); !equalIgnoringOrder(got, want) {
        t.Fatalf("subsets([0]) = %v, want %v", got, want)
    }
}
