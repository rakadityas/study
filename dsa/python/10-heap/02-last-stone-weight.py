# https://leetcode.com/problems/last-stone-weight/description/
# time complexity: O(nlogn) | n from the conversion, logn from the heap operations
# space complexity: O(n)

from typing import Optional, List
import heapq

class Solution:
    def lastStoneWeight(self, stones: list[int]) -> int:
        for i in range(len(stones)):
            stones[i] = stones[i] * -1
        
        heapq.heapify(stones)

        while len(stones) > 1:
            stoneOne = heapq.heappop(stones)
            stoneTwo = heapq.heappop(stones)

            if stoneOne < stoneTwo:
                heapq.heappush(stones, stoneOne - stoneTwo)
                    
        if len(stones) == 0:
            return 0 

        return stones[0] * -1
    
if "__name__" == "__main__":
    solution = Solution()
    assert solution.lastStoneWeight([2,7,4,1,8,1]) == 1
    assert solution.lastStoneWeight([1]) == 1
    assert solution.lastStoneWeight([1,1]) == 0
        
