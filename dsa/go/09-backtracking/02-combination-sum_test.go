package backtracking

import "testing"

// combinationSum returns every combination of candidates (reusable) that sums to target.
// Approach: backtrack from each start index, recursing on the same index so candidates can repeat,
// and pruning as soon as the running sum overshoots target
// time: exponential in target/min(candidates), space: O(target/min(candidates)) — recursion depth
func combinationSum(candidates []int, target int) [][]int {
    res := [][]int{}

    var backtracking func(start int, combination []int, currentSum int)
    backtracking = func(start int, combination []int, currentSum int) {
        if currentSum == target {
            res = append(res, append([]int{}, combination...))
            return
        }

        if currentSum > target { return }

        for i := start; i < len(candidates); i++ {
            combination = append(combination, candidates[i])
            backtracking(i, combination, currentSum+candidates[i])
            combination = combination[:len(combination)-1]
        }
    }

    backtracking(0, []int{}, 0)
    return res
}

func TestCombinationSum(t *testing.T) {
    cases := []struct {
        candidates []int
        target     int
        want       [][]int
    }{
        {[]int{2, 3, 6, 7}, 7, [][]int{{2, 2, 3}, {7}}},
        {[]int{2, 3, 5}, 8, [][]int{{2, 2, 2, 2}, {2, 3, 3}, {3, 5}}},
        {[]int{2}, 1, [][]int{}},
    }

    for _, c := range cases {
        if got := combinationSum(c.candidates, c.target); !equalIgnoringOrder(got, c.want) {
            t.Fatalf("combinationSum(%v, %d) = %v, want %v", c.candidates, c.target, got, c.want)
        }
    }
}
