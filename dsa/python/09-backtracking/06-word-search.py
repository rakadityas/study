# https://leetcode.com/problems/word-search/description/
# time complexity: O(m*n*4^len(word))
# space complexity: O(len(word))

from typing import List

class Solution:
    def exist(self, board: List[List[str]], word: str) -> bool:
        if not board:
            return False
        
        rows, cols = len(board), len(board[0])
        
        def backtracking(r: int, c: int, i: int) -> bool:
            if i == len(word):
                return True
            
            if r < 0 or r >= rows or c < 0 or c >= cols or board[r][c] != word[i]:
                return False
            
            temp = board[r][c]
            board[r][c] = "#"
            
            found = backtracking(r+1, c, i+1) or backtracking(r-1, c, i+1) or backtracking(r, c+1, i+1) or backtracking(r, c-1, i+1)
            
            board[r][c] = temp
            
            return found
        
        for r in range(rows):
            for c in range(cols):
                if backtracking(r, c, 0):
                    return True
        
        return False    

if __name__ == "__main__":
    board = [["A","B","C","E"],["S","F","C","S"],["A","D","E","E"]]
    word = "ABCCED"
    assert Solution().exist(board, word) == True
    
    board = [["A","B","C","E"],["S","F","C","S"],["A","D","E","E"]]
    word = "SEE"
    assert Solution().exist(board, word) == True
    
    board = [["A","B","C","E"],["S","F","C","S"],["A","D","E","E"]]
    word = "ABCB"
    assert Solution().exist(board, word) == False
