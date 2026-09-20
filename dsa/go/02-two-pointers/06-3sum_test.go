package two_pointers

import (
    "reflect"
    "sort"
    "testing"
)

// threeSum finds unique triplets that sum to zero.
// Approach: sort + two pointers; skip duplicate values in-place — no extra hashmap needed
// time: O(n²), space: O(1) extra — sort is O(n log n), dominated by the O(n²) outer+inner pointer scan
func threeSum(nums []int) [][]int {
    sort.Ints(nums)
    res := [][]int{}
    for i := 0; i < len(nums); i++ {
        if i > 0 && nums[i] == nums[i-1] { continue }
        l, r := i+1, len(nums)-1
        for l < r {
            s := nums[i] + nums[l] + nums[r]
            if s == 0 {
                res = append(res, []int{nums[i], nums[l], nums[r]})
                l++
                for l < r && nums[l] == nums[l-1] { l++ }
                r--
                for l < r && nums[r] == nums[r+1] { r-- }
            } else if s < 0 { l++ } else { r-- }
        }
    }
    return res
}

// threeSumMap finds the same triplets, deduplicating through a map instead of pointer skips.
// Approach: sort + two pointers, but let duplicates through and filter them with a map of seen triplets
// time: O(n²), space: O(n) — the map holds up to n result keys
func threeSumMap(nums []int) [][]int {
    nums = append([]int{}, nums...)
    sort.Ints(nums)
    mapResHistory := map[[3]int]bool{}
    res := [][]int{}

    for i := 0; i < len(nums); i++ {
        l, r := i+1, len(nums)-1

        for l < r {
            target := nums[i] + nums[r] + nums[l]
            if target == 0 {
                key := [3]int{nums[i], nums[l], nums[r]}
                if !mapResHistory[key] {
                    res = append(res, []int{nums[i], nums[l], nums[r]})
                    mapResHistory[key] = true
                }

                l++
                r--
            } else if target < 0 {
                l++
            } else {
                r--
            }
        }
    }

    return res
}

// threeSumSet finds the same triplets without sorting the input at all.
// Approach: fix i, then scan for a pair summing to -nums[i] with a hash set of values seen so far,
// deduplicating results through a set of sorted triplet keys
// time: O(n²), space: O(n) — one hash set per outer iteration plus the seen-triplet set
func threeSumSet(nums []int) [][]int {
    res := [][]int{}
    seen := map[[3]int]bool{}

    for i := 0; i < len(nums); i++ {
        target := -nums[i]
        twoSum := map[int]bool{}

        for j := i + 1; j < len(nums); j++ {
            complement := target - nums[j]

            if twoSum[complement] {
                triplet := []int{nums[i], nums[j], complement}
                sort.Ints(triplet)

                key := [3]int{triplet[0], triplet[1], triplet[2]}
                if !seen[key] {
                    seen[key] = true
                    res = append(res, triplet)
                }
            }

            twoSum[nums[j]] = true
        }
    }

    return res
}

// sortTriplets puts a result set in a canonical order so variants can be compared.
func sortTriplets(triplets [][]int) {
    sort.Slice(triplets, func(i, j int) bool {
        a, b := triplets[i], triplets[j]
        for k := 0; k < 3; k++ { if a[k] != b[k] { return a[k] < b[k] } }
        return false
    })
}

func TestThreeSum(t *testing.T) {
    cases := []struct {
        in   []int
        want [][]int
    }{
        {[]int{-1, 0, 1, 2, -1, -4}, [][]int{{-1, -1, 2}, {-1, 0, 1}}},
        {[]int{0, 1, 1}, [][]int{}},
        {[]int{0, 0, 0, 0}, [][]int{{0, 0, 0}}},
    }

    variants := []struct {
        name string
        fn   func([]int) [][]int
    }{
        {"threeSum", threeSum},
        {"threeSumMap", threeSumMap},
        {"threeSumSet", threeSumSet},
    }

    for _, c := range cases {
        want := append([][]int{}, c.want...)
        sortTriplets(want)

        for _, v := range variants {
            got := v.fn(append([]int{}, c.in...))
            sortTriplets(got)

            if len(got) == 0 && len(want) == 0 { continue }
            if !reflect.DeepEqual(got, want) {
                t.Fatalf("%s(%v) = %v, want %v", v.name, c.in, got, want)
            }
        }
    }
}

