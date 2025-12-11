# https://leetcode.com/problems/number-of-islands/description/
# time complexity: O(mn)
# space complexity: O(mn)

from typing import List

class Solution:
    def numIslands(self, grid: List[List[str]]) -> int:
        directions = [[1,0],[-1,0],[0,1],[0,-1]]
        maxRows, maxCols = len(grid), len(grid[0])
        islands = 0

        def dfs(row: int, col: int):
            if row < 0 or col < 0 or row >= maxRows or col >= maxCols or grid[row][col] == "0":
                return 

            grid[row][col] = "0"
            for dRow, dCol in directions:
                dfs(row + dRow, col + dCol)

        for r in range(len(grid)):
            for c in range(len(grid[r])):
                if grid[r][c] == "1":
                    dfs(r, c)
                    islands += 1
        
        return islands

if __name__ == "__main__":
    s = Solution()
    assert s.numIslands([["1","1","1","1","0"],["1","1","0","1","0"],["1","1","0","0","0"],["0","0","0","0","0"]]) == 1
    assert s.numIslands([["1","1","0","0","0"],["1","1","0","0","0"],["0","0","1","0","0"],["0","0","0","1","1"]]) == 3