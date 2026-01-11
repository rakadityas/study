# https://leetcode.com/problems/longest-consecutive-sequence/description/
# time complexity: O(n)
# space complexity: O(n)

from collections import defaultdict

class Solution:
    def longestConsecutive(self, nums: list[int]) -> int:
        if not nums:
            return 0

        mapHistory = defaultdict(bool)
        for n in nums:
            mapHistory[n] = True

        maxSequence = 0

        for n in mapHistory:
            if n - 1 in mapHistory:
                continue

            counter = 0
            while n + counter in mapHistory:
                counter += 1
            
            maxSequence = max(maxSequence, counter)

        return maxSequence
    

if __name__ == "__main__":
    solution = Solution()

    assert solution.longestConsecutive([100,4,200,1,3,2]) == 4
    assert solution.longestConsecutive([0,3,7,2,5,8,4,6,0,1]) == 9
    assert solution.longestConsecutive([1,0,1,2]) == 3
    assert solution.longestConsecutive([]) == 0