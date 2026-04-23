# https://leetcode.com/problems/kth-largest-element-in-an-array/description/

import heapq

# Approach: negate all values → heapify the whole array into a max-heap → pop k-1 times to reach kth largest
# time complexity: O(n log n) — heapify is O(n), but k pops are each O(log n); total O(n log n) worst case
# space complexity: O(n) — heapify modifies in-place but logically uses O(n) heap space
class Solution:
    def findKthLargest(self, nums: list[int], k: int) -> int:
        for i in range(len(nums)):
            nums[i] = nums[i] * -1

        heapq.heapify(nums)

        while k > 1:
            heapq.heappop(nums)
            k = k - 1
        
        return heapq.heappop(nums) * -1

# Approach: maintain a min-heap of size k; push each element and pop if heap exceeds k → root is the kth largest
# time complexity: O(n log k) — each of the n elements does a push+conditional pop into a heap of size k
# space complexity: O(k) — heap never exceeds k elements
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