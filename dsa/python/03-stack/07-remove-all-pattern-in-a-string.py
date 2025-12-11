# not found
# time complexity: O(n)
# space complexity: O(n)

class Solution:
    def removeOccurrences(self, s: str) -> str:
        targetDict = {"ab", "bc", "cd", "dc", "cb", "ba"}
        stack = []

        for i in range(len(s)):
            if stack and stack[-1] + s[i] in targetDict:
                stack.pop()
            else:
                stack.append(s[i])
            
        return "".join(stack)

if __name__ == "__main__":
    s = Solution()
    assert s.removeOccurrences("abccba") == "ca"
    assert s.removeOccurrences("aabccba") == "aca"
    assert s.removeOccurrences("") == ""