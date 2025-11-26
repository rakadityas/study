# https://leetcode.com/problems/3sum/description/

from typing import List

# time complexity: O(n^2)
# space complexity: O(n)
class SolutionOne:
    def threeSum(self, nums: List[int]) -> List[List[int]]:
        lenNums = len(nums)-1
        nums.sort()
        mapResHistory = {}
        res = []

        for i in range(len(nums)):
            l, r = i+1, lenNums

            while l < r:
                target = nums[i] + nums[r] + nums[l]
                if target == 0:
                    tupleNums = tuple([nums[i], nums[r], nums[l]])
                    if tupleNums not in mapResHistory:
                        res.append([nums[i], nums[r], nums[l]])
                        mapResHistory[tupleNums] = True
                    
                    l += 1
                    r -= 1

                elif target < 0:
                    l += 1
                else:
                    r -= 1        
        return res

# time complexity: O(n^2)
# space complexity: O(1)
class SolutionTwo:
    def threeSum(self, nums: List[int]) -> List[List[int]]:
        nums.sort()
        res = []
        n = len(nums)

        for i in range(n):
            if i > 0 and nums[i] == nums[i - 1]:
                continue  # skip duplicates for i

            l, r = i + 1, n - 1
            
            while l < r:
                s = nums[i] + nums[l] + nums[r]

                if s < 0:
                    l += 1
                elif s > 0:
                    r -= 1
                else:
                    res.append([nums[i], nums[l], nums[r]])
                    l += 1
                    r -= 1

                    # skip duplicates for left pointer
                    while l < r and nums[l] == nums[l - 1]:
                        l += 1

                    # skip duplicates for right pointer
                    while l < r and nums[r] == nums[r + 1]:
                        r -= 1
        
        return res

if __name__ == "__main__":
    solution = SolutionOne()
    assert solution.threeSum([-1,0,1,2,-1,-4]) == [[-1, 2, -1], [-1, 1, 0]]
    assert solution.threeSum([0,1,1]) == []
    assert solution.threeSum([0,0,0]) == [[0,0,0]]
    assert solution.threeSum([0,0,0,0]) == [[0,0,0]]

    solution = SolutionTwo()
    assert solution.threeSum([-1,0,1,2,-1,-4]) == [[-1, -1, 2], [-1, 0, 1]]
    assert solution.threeSum([0,1,1]) == []
    assert solution.threeSum([0,0,0]) == [[0,0,0]]
    assert solution.threeSum([0,0,0,0]) == [[0,0,0]]
        
