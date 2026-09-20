package backtracking

import (
    "sort"
    "strconv"
    "strings"
    "testing"
)

// combinationSum2Map returns every distinct combination using each candidate at most once.
// Approach: sort descending, backtrack from i+1, and deduplicate results via a map of combination keys
// time: O(2^n), space: O(n) for recursion plus O(results) for the map — key building adds a per-result cost
func combinationSum2Map(candidates []int, target int) [][]int {
    candidates = append([]int{}, candidates...)
    sort.Sort(sort.Reverse(sort.IntSlice(candidates)))
    res := [][]int{}
    mapHistory := map[string]bool{}

    key := func(combination []int) string {
        parts := make([]string, 0, len(combination))
        for _, v := range combination {
            parts = append(parts, strconv.Itoa(v))
        }
        return strings.Join(parts, ",")
    }

    var backtracking func(combination []int, currentSum, start int)
    backtracking = func(combination []int, currentSum, start int) {
        if currentSum == target {
            k := key(combination)
            if !mapHistory[k] {
                res = append(res, append([]int{}, combination...))
                mapHistory[k] = true
            }
            return
        } else if currentSum > target {
            return
        }

        for i := start; i < len(candidates); i++ {
            combination = append(combination, candidates[i])
            backtracking(combination, currentSum+candidates[i], i+1)
            combination = combination[:len(combination)-1]
        }
    }

    backtracking([]int{}, 0, 0)
    return res
}

// combinationSum2Skip returns the same combinations without a dedup map.
// Approach: sort ascending, then at each level skip a candidate equal to the previous one tried,
// which prevents two branches from producing the same combination
// time: O(2^n), space: O(n) — recursion stack only
func combinationSum2Skip(candidates []int, target int) [][]int {
    candidates = append([]int{}, candidates...)
    sort.Ints(candidates)
    res := [][]int{}

    var backtracking func(path []int, currentSum, start int)
    backtracking = func(path []int, currentSum, start int) {
        if currentSum == target {
            res = append(res, append([]int{}, path...))
            return
        } else if currentSum > target {
            return
        }

        prev := 0
        hasPrev := false
        for i := start; i < len(candidates); i++ {
            if hasPrev && candidates[i] == prev { continue }

            path = append(path, candidates[i])
            backtracking(path, currentSum+candidates[i], i+1)
            path = path[:len(path)-1]

            prev, hasPrev = candidates[i], true
        }
    }

    backtracking([]int{}, 0, 0)
    return res
}

func TestCombinationSum2(t *testing.T) {
    cases := []struct {
        candidates []int
        target     int
        want       [][]int
    }{
        {[]int{10, 1, 2, 7, 6, 1, 5}, 8, [][]int{{1, 1, 6}, {1, 2, 5}, {1, 7}, {2, 6}}},
        {[]int{2, 5, 2, 1, 2}, 5, [][]int{{1, 2, 2}, {5}}},
    }

    for _, c := range cases {
        // the map variant sorts descending, so normalize each combination before comparing
        got := combinationSum2Map(c.candidates, c.target)
        for i := range got { sort.Ints(got[i]) }
        if !equalIgnoringOrder(got, c.want) {
            t.Fatalf("combinationSum2Map(%v, %d) = %v, want %v", c.candidates, c.target, got, c.want)
        }

        if got := combinationSum2Skip(c.candidates, c.target); !equalIgnoringOrder(got, c.want) {
            t.Fatalf("combinationSum2Skip(%v, %d) = %v, want %v", c.candidates, c.target, got, c.want)
        }
    }
}
