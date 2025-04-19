from collections import defaultdict

class Solution:
    def longestConsecutive(self, nums: list[int]) -> int:

        maxSequence = 0

        mapHistory = defaultdict(bool)
        for i in range(len(nums)):
            mapHistory[nums[i]] = True

        for i in range(len(nums)):
            if nums[i]-1 in mapHistory:
                continue

            counter = 0
            while nums[i]+counter in mapHistory:
                counter += 1
            
            if counter > maxSequence:
                maxSequence = counter
        
        return maxSequence
    

if __name__ == "__main__":
    solution = Solution()

    assert solution.longestConsecutive([100,4,200,1,3,2]) == 4
    assert solution.longestConsecutive([0,3,7,2,5,8,4,6,0,1]) == 9
    assert solution.longestConsecutive([1,0,1,2]) == 3
    assert solution.longestConsecutive([]) == 0