// https://leetcode.com/problems/contains-duplicate/description/

package array_hashes

import (
	"testing"
)

func ContainsDuplicate(nums []int) bool {
	mapHistory := make(map[int]bool)

	for i := range nums {
		if mapHistory[nums[i]] {
			return true
		}

		mapHistory[nums[i]] = true
	}

	return false
}

func TestContainsDuplicate(t *testing.T) {
	testCases := []struct {
		name     string
		nums     []int
		expected bool
	}{
		{
			name:     "contains duplicate",
			nums:     []int{1, 2, 3, 1},
			expected: true,
		},
		{
			name:     "no duplicate",
			nums:     []int{1, 2, 3, 4},
			expected: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := ContainsDuplicate(tc.nums)
			if result != tc.expected {
				t.Errorf("ContainsDuplicate(%v) = %v; want %v", tc.nums, result, tc.expected)
			}
		})
	}
}
