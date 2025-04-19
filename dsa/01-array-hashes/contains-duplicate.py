# https://leetcode.com/problems/contains-duplicate/description/

from collections import defaultdict

class Solution:
    def containsDuplicate(self, nums: list[int]) -> bool:
        mapHistory = defaultdict(bool)

        for i in range(len(nums)):
            if nums[i] in mapHistory:
                return True
            
            mapHistory[nums[i]] = True
        
        return False
    
if __name__ == '__main__':
    solution = Solution()

    assert solution.containsDuplicate([1,2,3,1]) ==  True
    assert solution.containsDuplicate([1,2,3,4]) == False
    assert solution.containsDuplicate([1,1,1,3,3,4,3,2,4,2]) == True
