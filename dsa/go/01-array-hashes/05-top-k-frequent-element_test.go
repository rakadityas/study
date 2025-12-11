package array_hashes

import (
    "reflect"
    "sort"
    "testing"
)

// topKFrequent returns the k most frequent elements in nums.
// Uses bucket sort on frequency counts for O(n) average time.
func topKFrequent(nums []int, k int) []int {
    if k == 0 || len(nums) == 0 {
        return []int{}
    }

    // Count frequencies
    freq := make(map[int]int)
    for _, n := range nums {
        freq[n]++
    }

    // Buckets where index = frequency, value = list of numbers with that frequency
    buckets := make([][]int, len(nums)+1)
    for num, count := range freq {
        buckets[count] = append(buckets[count], num)
    }

    // Collect top k from highest frequency bucket down
    res := make([]int, 0, k)
    for i := len(buckets) - 1; i >= 0 && len(res) < k; i-- {
        for _, num := range buckets[i] {
            res = append(res, num)
            if len(res) == k {
                break
            }
        }
    }

    return res
}

func TestTopKFrequent(t *testing.T) {
    testCases := []struct {
        name     string
        nums     []int
        k        int
        expected []int
    }{
        {
            name:     "example 1",
            nums:     []int{1, 1, 1, 2, 2, 3},
            k:        2,
            expected: []int{1, 2},
        },
        {
            name:     "example 2",
            nums:     []int{1},
            k:        1,
            expected: []int{1},
        },
        {
            name:     "ties and multiple",
            nums:     []int{4, 4, 4, 5, 5, 6, 6, 6, 7},
            k:        3,
            expected: []int{4, 6, 5},
        },
    }

    for _, tc := range testCases {
        t.Run(tc.name, func(t *testing.T) {
            got := topKFrequent(tc.nums, tc.k)
            // Sort both for comparison since order among same frequency is not guaranteed
            sort.Ints(got)
            exp := make([]int, len(tc.expected))
            copy(exp, tc.expected)
            sort.Ints(exp)

            if !reflect.DeepEqual(got, exp) {
                t.Errorf("topKFrequent(%v, %d) = %v; want %v", tc.nums, tc.k, got, exp)
            }
        })
    }
}

