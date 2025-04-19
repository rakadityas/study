from collections import defaultdict

class Solution:
    def twoSum(self, nums: list[int], target: int) -> list[int]:
        mapHistory = defaultdict(int)

        for i in range(len(nums)):
            gapVal = target-nums[i]
            if gapVal in mapHistory:
                return [mapHistory[gapVal], i]
            
            mapHistory[nums[i]] = i
        
        return [0,0]

if __name__ == "__main__":
    solution = Solution()

    assert solution.twoSum([2,7,11,15], 9) == [0,1]
    assert solution.twoSum([3,2,4], 6) == [1,2]
    assert solution.twoSum([3,3], 6) == [0,1]
    assert solution.twoSum([2,7,11,15], 100) == [0,0]
