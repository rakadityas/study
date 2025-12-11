# https://leetcode.com/problems/remove-element/description/
# time complexity: O(n)
# space complexity: O(1)

from typing import List

class Solution:
    def removeElement(self, nums: List[int], val: int) -> int:
        l = 0

        for i in range(len(nums)):
            if nums[i] != val:
                nums[l] = nums[i]
                l += 1
        return l

if __name__ == "__main__":
    solution = Solution()
    assert solution.removeElement([3,2,2,3], 3) == 2
    assert solution.removeElement([0,1,2,2,3,0,4,2], 2) == 5