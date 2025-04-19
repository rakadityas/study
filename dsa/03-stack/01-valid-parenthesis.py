
class Solution:
    def isValid(self, s:str) -> bool:
        
        mapBracket = {
            "}": "{",
            ")": "(",
            "]": "["
        }

        stack = []

        for i in range(len(s)):
            if s[i] not in mapBracket:
                stack.append(s[i])
                continue

            if not stack:
                return False
            
            val = stack.pop()
            if mapBracket[s[i]] == val:
                continue
            else:
                return False
        
        return len(stack) == 0
    
if __name__ == "__main__":
    solution = Solution()

    assert solution.isValid('()') == True
    assert solution.isValid('()[]{}') == True
    assert solution.isValid('(]') == False
    assert solution.isValid('[]') == True

