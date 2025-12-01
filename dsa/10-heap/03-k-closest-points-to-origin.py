# https://leetcode.com/problems/k-closest-points-to-origin/description/
# time complexity: O(nlogk)
# space complexity: O(k)

from typing import List
import heapq

class Solution:
    def kClosest(self, points: List[List[int]], k: int) -> List[List[int]]:
        heapVal = []
        heapq.heapify(heapVal)

        for i in range(len(points)):
            calcRes = (points[i][0] ** 2) + (points[i][1] ** 2)
            heapq.heappush(heapVal, (calcRes, points[i]))
        
        outputRes = []

        while k > 0:
            outputRes.append(heapq.heappop(heapVal)[1])
            k = k -1
        
        return outputRes

if "__name__" == "__main__":
    solution = Solution()
    assert solution.kClosest([[1,3],[-2,2]], 1) == [[-2,2]]
    assert solution.kClosest([[3,3],[5,-1],[-2,4]], 2) == [[3,3],[-2,4]]
