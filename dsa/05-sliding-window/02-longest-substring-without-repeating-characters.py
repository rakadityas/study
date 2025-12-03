# https://leetcode.com/problems/longest-substring-without-repeating-characters/description/
# time complexity: O(n)
# space complexity: O(n)

from collections import defaultdict

class Solution:
    def lenghthOfLongestSubstring(self, s:str) -> int:
        left = 0

        maxCounter = 0
        
        mapHistory = defaultdict(bool)

        for i in range(len(s)):
            while s[i] in mapHistory:
                del mapHistory[s[left]]
                left += 1
            
            mapHistory[s[i]] = True
    
            maxCounter = max(maxCounter, i-left+1)

        return maxCounter
    
    def lengthOfLongestSubstringLastSeen(self, s: str) -> int:
        left = 0 
        mapHistory = defaultdict(int)
        maxCounter = 0
        
        for i in range(len(s)):
            if s[i] in mapHistory:
                left = max(left, mapHistory[s[i]]+1)
            mapHistory[s[i]] = i
            maxCounter = max(maxCounter, i-left+1)
        return maxCounter
        

if __name__ == "__main__":
    solution = Solution()

    assert solution.lenghthOfLongestSubstring("abcabcbb") == 3
    assert solution.lenghthOfLongestSubstring("bbbbb") == 1
    assert solution.lenghthOfLongestSubstring("pwwkew") == 3

    assert solution.lengthOfLongestSubstringLastSeen("abcabcbb") == 3
    assert solution.lengthOfLongestSubstringLastSeen("bbbbb") == 1
    assert solution.lengthOfLongestSubstringLastSeen("pwwkew") == 3

