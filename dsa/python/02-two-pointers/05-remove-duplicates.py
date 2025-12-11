# https://leetcode.com/problems/remove-duplicates-from-sorted-array/description/
# time complexity: O(n)
# space complexity: O(1)

from typing import List
from collections import defaultdict

class Solution:
    def removeDuplicates(self, nums: List[int]) -> int:
        mapHistory = defaultdict(int)

        l = 0

        for i in range(len(nums)):
            if nums[i] in mapHistory:
                continue
            
            mapHistory[nums[i]] = True
            nums[l] = nums[i]
            l += 1
        
        return l
    
if __name__ == "__main__":
    solution = Solution()
    assert solution.removeDuplicates([1,1,2]) == 2
    assert solution.removeDuplicates([0,0,1,1,1,2,2,3,3,4]) == 5