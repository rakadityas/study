# https://leetcode.com/problems/k-closest-points-to-origin/description/
# time complexity: O(nlogk)
# space complexity: O(k)

from typing import List
import heapq

class Solution:
    def kClosest(self, points: List[List[int]], k: int) -> List[List[int]]:
        heap = []

        for i in range(len(points)):
            dist = points[i][0]*points[i][0] + points[i][1]*points[i][1]
            heapq.heappush(heap, (dist * -1, points[i]))

            if len(heap) > k:
                heapq.heappop(heap)

        result = []
        for i in range(len(heap)):
            result.append(heap[i][1])
        
        return result

if "__name__" == "__main__":
    solution = Solution()
    assert solution.kClosest([[1,3],[-2,2]], 1) == [[-2,2]]
    assert solution.kClosest([[3,3],[5,-1],[-2,4]], 2) == [[3,3],[-2,4]]
