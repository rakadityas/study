# https://leetcode.com/problems/remove-all-adjacent-duplicates-in-string/description/
# time complexity: O(n)
# space complexity: O(n)

class Solution:
    def removeDuplicates(self, s: str) -> str:
        stack = []

        for i in range(len(s)):
            if stack and stack[-1] == s[i]:
                stack.pop()
            else:
                stack.append(s[i])
            
        return "".join(stack)

if __name__ == "__main__":
    s = Solution()
    assert s.removeDuplicates("abbaca") == "ca"
    assert s.removeDuplicates("azxxzy") == "ay"
