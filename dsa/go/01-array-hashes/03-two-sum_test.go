package array_hashes

import (
	"reflect"
	"testing"
)

func twoSum(nums []int, target int) []int {
	mapHistory := make(map[int]int)

	for i := range nums {
		gapVal := target - nums[i]
		if _, ok := mapHistory[gapVal]; ok {
			return []int{mapHistory[gapVal], i}
		}
		mapHistory[nums[i]] = i
	}

	return []int{0, 0}
}

func TestTwoSum(t *testing.T) {
	testCases := []struct {
		name     string
		nums     []int
		target   int
		expected []int
	}{
		{
			name:     "two sum",
			nums:     []int{2, 7, 11, 15},
			target:   9,
			expected: []int{0, 1},
		},
		{
			name:     "two sum",
			nums:     []int{3, 2, 4},
			target:   6,
			expected: []int{1, 2},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := twoSum(tc.nums, tc.target)
			if !reflect.DeepEqual(result, tc.expected) {
				t.Errorf("twoSum(%v, %v) = %v; want %v", tc.nums, tc.target, result, tc.expected)
			}
		})
	}
}
