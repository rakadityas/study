# https://leetcode.com/problems/remove-all-occurrences-of-a-substring/description/
# time complexity: O(n * m)
# space complexity: O(n)

class Solution:
    def removeOccurrences(self, s: str, part: str) -> str:
        stack = []
        lenPart = len(part)

        for i in range(len(s)):
            stack.append(s[i])

            if len(stack) >= len(part) and "".join(stack[lenPart*-1:]) == part:
                for j in range(lenPart):
                    stack.pop()
        
        return "".join(stack)

if __name__ == "__main__":
    solution = Solution()
    
    assert solution.removeOccurrences("daabcbaabcbc", "abc") == "dab"
    assert solution.removeOccurrences("axxxxyyyyb", "xy") == "ab"