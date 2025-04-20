from collections import defaultdict

class Solution:
    def lenghthOfLongestSubstring(self, s:str) -> int:
        l = 0
        r = 0

        maxCounter = 0
        
        mapHistory = defaultdict(bool)

        while r < len(s):
            while s[r] in mapHistory:
                del mapHistory[s[l]]
                l += 1
            
            mapHistory[s[r]] = True
            maxCounter = max(maxCounter, r-l+1)
            
            r += 1

        return maxCounter

if __name__ == "__main__":
    solution = Solution()

    print(solution.lenghthOfLongestSubstring("abcabcbb"))

    assert solution.lenghthOfLongestSubstring("abcabcbb") == 3
    assert solution.lenghthOfLongestSubstring("bbbbb") == 1
    assert solution.lenghthOfLongestSubstring("pwwkew") == 3

