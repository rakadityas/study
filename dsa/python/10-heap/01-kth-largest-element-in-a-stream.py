# https://leetcode.com/problems/kth-largest-element-in-a-stream/description/

from typing import List
import heapq

# Approach: min-heap of size k; after each push, pop until heap size == k via a separate maintainHeap loop
# time complexity: O(n log k) — each add does a push + conditional pops
# space complexity: O(k) — heap is capped at k elements
class KthLargest:
    def __init__(self, k: int, nums: List[int]):
        self.minHeap = nums
        self.k = k

        heapq.heapify(self.minHeap)
        self.maintainHeap() 
        return

    def add(self, val: int) -> int:
        heapq.heappush(self.minHeap, val)
        self.maintainHeap()
        return self.minHeap[0]

    def maintainHeap(self):
        while len(self.minHeap) > self.k:
            heapq.heappop(self.minHeap)
        return

# Approach: same min-heap of size k, but uses heapreplace (atomic pop+push) to avoid unnecessary heap resizes
# time complexity: O(n log k) — same asymptotic, but heapreplace is faster in practice (one heap operation vs two)
# space complexity: O(k)
class KthLargestOptimal:
    def __init__(self, k: int, nums: List[int]):
        self.k = k
        self.minHeap = nums

        heapq.heapify(self.minHeap)

        while len(self.minHeap) > k:
            heapq.heappop(self.minHeap)

    def add(self, val: int) -> int:
        if len(self.minHeap) < self.k:
            heapq.heappush(self.minHeap, val)
        elif val > self.minHeap[0]:
            heapq.heapreplace(self.minHeap, val)

        return self.minHeap[0]

if __name__ == "__main__":
    kthLargest = KthLargest(3, [4, 5, 8, 2])
    assert kthLargest.add(3) == 4
    assert kthLargest.add(5) == 5
    assert kthLargest.add(10) == 5
    assert kthLargest.add(9) == 8
    assert kthLargest.add(4) == 8
    
    kthLargest = KthLargest(4, [7,7,7,7,8,3])
    assert kthLargest.add(2) == 7
    assert kthLargest.add(10) == 7
    assert kthLargest.add(9) == 7
    assert kthLargest.add(9) == 8
    