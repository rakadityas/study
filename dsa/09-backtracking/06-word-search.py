# https://leetcode.com/problems/word-search/

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
class SolutionBrute: # not working really well, unfinished, and only covers the right -> bottom -> left direction
    def exist(self, board: List[List[str]], word: str) -> bool:
        # right -> bottom -> left
        lenWord = len(word)
        self.hasFound = False

        def backtrack(path: List[str], wordIdx: int, startI: int, startJ: int):
            print(path)
            if board[startI][startJ] != word[wordIdx]:
                return
            else:
                wordIdx += 1

            if lenWord == wordIdx:
                self.hasFound = True

            for i in range(startI, len(board)):
                for j in range(startJ, len(board)):
                    path.append(board[i][j])
                    backtrack(path, wordIdx+1, i+1, j+1) #traverse right
                    path.pop()

                    if self.hasFound:
                        break
                    
                    for iBottom in range(i, len(board)):
                        path.append(board[iBottom][j])
                        backtrack(path, wordIdx+1, iBottom+1, j) #traverse bottom
                        path.pop()
                    
                        if self.hasFound:
                            break
                    
                    for iLeft in range(i, -1, -1):
                        path.append(board[iLeft][j])
                        backtrack(path, wordIdx+1, iLeft, j-1) #traverse left
                        path.pop()

                        if self.hasFound:
                            break
                
                if self.hasFound:
                        break
            return
        
        backtrack([], 0, 0, 0)

        return self.hasFound
    

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
