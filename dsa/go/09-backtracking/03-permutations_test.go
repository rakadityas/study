package backtracking

import "testing"

// permuteList returns every permutation of nums.
// Approach: boolean record array to track used indices — faster than hashing, direct index access
// time: O(n*n!), space: O(n) — n! permutations each costing O(n) to copy; record array + recursion stack
func permuteList(nums []int) [][]int {
    res := [][]int{}
    record := make([]bool, len(nums))

    var backtracking func(permutation []int)
    backtracking = func(permutation []int) {
        if len(permutation) == len(nums) {
            res = append(res, append([]int{}, permutation...))
            return
        }

        for i := range nums {
            if record[i] { continue }

            record[i] = true
            permutation = append(permutation, nums[i])
            backtracking(permutation)
            permutation = permutation[:len(permutation)-1]
            record[i] = false
        }
    }

    backtracking([]int{})
    return res
}

// permuteMap returns the same permutations using a map of used indices.
// Approach: map to track used indices — O(1) average lookup, but with hashing overhead a plain array avoids
// time: O(n*n!), space: O(n) — map size bounded by n plus the recursion stack
func permuteMap(nums []int) [][]int {
    res := [][]int{}
    history := map[int]bool{}

    var backtracking func(permutation []int)
    backtracking = func(permutation []int) {
        if len(permutation) == len(nums) {
            res = append(res, append([]int{}, permutation...))
            return
        }

        for i := range nums {
            if history[i] { continue }

            history[i] = true
            permutation = append(permutation, nums[i])
            backtracking(permutation)
            permutation = permutation[:len(permutation)-1]
            delete(history, i)
        }
    }

    backtracking([]int{})
    return res
}

func TestPermute(t *testing.T) {
    cases := []struct {
        in   []int
        want [][]int
    }{
        {[]int{1, 2, 3}, [][]int{{1, 2, 3}, {1, 3, 2}, {2, 1, 3}, {2, 3, 1}, {3, 1, 2}, {3, 2, 1}}},
        {[]int{0, 1}, [][]int{{0, 1}, {1, 0}}},
        {[]int{1}, [][]int{{1}}},
    }

    for _, c := range cases {
        if got := permuteList(c.in); !equalIgnoringOrder(got, c.want) {
            t.Fatalf("permuteList(%v) = %v, want %v", c.in, got, c.want)
        }
        if got := permuteMap(c.in); !equalIgnoringOrder(got, c.want) {
            t.Fatalf("permuteMap(%v) = %v, want %v", c.in, got, c.want)
        }
    }
}
