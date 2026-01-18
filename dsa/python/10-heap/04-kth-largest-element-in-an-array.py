# https://leetcode.com/problems/kth-largest-element-in-an-array/description/

import heapq

# time complexity: O(nlogn): n for the conversion, logn for the heap operations, because we are converting everything first then it will be nlogn
# space complexity: O(n)
class Solution:
    def findKthLargest(self, nums: list[int], k: int) -> int:
        for i in range(len(nums)):
            nums[i] = nums[i] * -1

        heapq.heapify(nums)

        while k > 1:
            heapq.heappop(nums)
            k = k - 1
        
        return heapq.heappop(nums) * -1

# time complexity: O(nlogk): k is the size of the heap, n is the number of elements in the array
# space complexity: O(k)
class SolutionOptimal:
    def findKthLargest(self, nums: List[int], k: int) -> int:
        heap = []
        heapq.heapify(heap)

        for i in range(len(nums)):
            heapq.heappush(heap, nums[i])
            if len(heap) > k:
                heapq.heappop(heap)
        
        return heap[0]

if __name__ == "__main__":
    s = Solution()
    assert s.findKthLargest([3,2,1,5,6,4], 2) == 5
    assert s.findKthLargest([3,2,3,1,2,4,5,5,6], 4) == 4