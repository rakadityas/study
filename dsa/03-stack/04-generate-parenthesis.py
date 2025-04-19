
class Solution:
    def __init__(self):
        self.stack = []
        self.response = []

    def generateParenthesis(self, n: int) -> list[str]:
        self.generate(n, 0, 0)
        return self.response
        
    def generate(self, n: int, open: int, close: int):
        if open == close == n:
            copyStack = self.stack.copy()
            self.response.append("".join(copyStack))
        
        if open < n:
            self.stack.append("(")
            self.generate( n, open+1, close)
            self.stack.pop()
        
        if close < open:
            self.stack.append(")")
            self.generate(n, open, close+1)
            self.stack.pop()

if __name__ == "__main__":
    solution = Solution()

    assert solution.generateParenthesis(3) == ["((()))", "(()())", "(())()", "()(())", "()()()"]
    assert solution.generateParenthesis(1) == ["()"]