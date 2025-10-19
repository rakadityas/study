# https://leetcode.com/problems/kth-largest-element-in-an-array/description/
# time complexity: O(nlogk)
# space complexity: O(k)

import heapq

class Solution:
    def findKthLargest(self, nums: list[int], k: int) -> int:
        for i in range(len(nums)):
            nums[i] = nums[i] * -1

        heapq.heapify(nums)

        while k > 1:
            heapq.heappop(nums)
            k = k - 1
        
        return heapq.heappop(nums) * -1

if __name__ == "__main__":
    s = Solution()
    assert s.findKthLargest([3,2,1,5,6,4], 2) == 5
    assert s.findKthLargest([3,2,3,1,2,4,5,5,6], 4) == 4