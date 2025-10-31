# https://leetcode.com/problems/max-area-of-island/description/
# time complexity: O(mn)
# space complexity: O(mn)

from typing import List

class Solution:
    def maxAreaOfIsland(self, grid: List[List[int]]) -> int:
        directions = [(0,1),(1,0),(-1,0),(0,-1)]
        maxRow, maxCol = len(grid), len(grid[0])
        maxArea = 0

        def dfs(row: int, col: int) -> int:
            if row < 0 or col < 0 or row >= maxRow or col >= maxCol or grid[row][col] == 0:
                return 0
            
            grid[row][col] = 0
            
            area = 1
            for dRow, dCol in directions:
                area += dfs(row + dRow, col + dCol)
            
            return area
        
        for r in range(maxRow):
            for c in range(maxCol):
                if grid[r][c] == 1:
                    area = dfs(r, c)
                    maxArea = max(maxArea, area)
        
        return maxArea