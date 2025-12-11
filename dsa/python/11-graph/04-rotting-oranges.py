from collections import deque
from typing import List

class Solution:
    def orangesRotting(self, grid: List[List[int]]) -> int:
        directionArr = [[0,1],[1,0],[-1,0],[0,-1]]
        maxRow, maxCol = len(grid), len(grid[0])
        rotQ = deque()
        fresh = 0
        minutes = 0

        for r in range(len(grid)):
            for c in range(len(grid[r])):
                if grid[r][c] == 2:
                    rotQ.append([r,c])
                if grid[r][c] == 1:
                    fresh += 1
        
        while rotQ and fresh > 0:
            lenRotQ = len(rotQ)
            
            for i in range(len(rotQ)):
                rotRow, rotCol = rotQ.popleft()

                for dirR, dirC in directionArr:
                    currR, currC = rotRow+dirR, rotCol+dirC
                    if currR >= 0 and currR < maxRow and currC >= 0 and currC < maxCol and grid[currR][currC] == 1:
                        grid[currR][currC] = 2
                        rotQ.append([currR, currC])
                        fresh -= 1

            minutes += 1
            
        if fresh == 0:
            return minutes

        return -1


if __name__ == "__main__":
    s = Solution()
    assert s.orangesRotting([[2,1,1],[1,1,0],[0,1,1]]) == 4
    assert s.orangesRotting([[2,1,1],[0,1,1],[1,0,1]]) == -1
    assert s.orangesRotting([[0,2]]) == 0