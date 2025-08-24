from typing import List
import heapq

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
    