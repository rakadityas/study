package array_hashes

import (
    "reflect"
    "sort"
    "testing"
)

// topKFrequent returns the k most frequent elements in nums.
// Approach: bucket sort — index bucket[freq] holds all numbers at that frequency, scan high→low
// time: O(n), space: O(n) — no comparison-based sort; frequency range is bounded by n
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

// topKFrequentMaxScan returns the same result by repeatedly taking the highest remaining frequency.
// Approach: count frequencies, group numbers by frequency, then pull whole groups off the top by
// scanning for the max key until k values have been collected
// time: O(n*k) worst case, space: O(n) — the max scan over distinct frequencies runs up to k times
func topKFrequentMaxScan(nums []int, k int) []int {
    mapNumsHistory := map[int]int{}
    for _, n := range nums {
        mapNumsHistory[n]++
    }

    mapGroupNums := map[int][]int{}
    for key, value := range mapNumsHistory {
        mapGroupNums[value] = append(mapGroupNums[value], key)
    }

    res := []int{}
    for k > 0 && len(mapGroupNums) > 0 {
        maxKey := 0
        for freq := range mapGroupNums {
            if freq > maxKey { maxKey = freq }
        }

        for _, val := range mapGroupNums[maxKey] {
            res = append(res, val)
            k--
            if k == 0 { return res }
        }

        delete(mapGroupNums, maxKey)
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
            // Sort both for comparison since order among same frequency is not guaranteed
            exp := make([]int, len(tc.expected))
            copy(exp, tc.expected)
            sort.Ints(exp)

            got := topKFrequent(tc.nums, tc.k)
            sort.Ints(got)
            if !reflect.DeepEqual(got, exp) {
                t.Errorf("topKFrequent(%v, %d) = %v; want %v", tc.nums, tc.k, got, exp)
            }

            gotMaxScan := topKFrequentMaxScan(tc.nums, tc.k)
            sort.Ints(gotMaxScan)
            if !reflect.DeepEqual(gotMaxScan, exp) {
                t.Errorf("topKFrequentMaxScan(%v, %d) = %v; want %v", tc.nums, tc.k, gotMaxScan, exp)
            }
        })
    }
}

