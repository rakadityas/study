# https://leetcode.com/problems/permutation-in-string/description/
# time complexity: O(m+n)
# space complexity: O(1)

class Solution:
    def checkInclusion(self, s1: str, s2: str) -> bool:
        if len(s1) > len(s2):
            return False
        
        freqS1 = [0] * 26
        freqWindow = [0] * 26
        
        for ch in s1:
            freqS1[ord(ch) - ord('a')] += 1
        
        windowSize = len(s1)
        left = 0
        
        for right in range(len(s2)):
            freqWindow[ord(s2[right]) - ord('a')] += 1
            
            if right - left + 1 > windowSize:
                freqWindow[ord(s2[left]) - ord('a')] -= 1
                left += 1
            
            if right - left + 1 == windowSize and freqWindow == freqS1:
                return True
        
        return False

if __name__ == "__main__":
    s = Solution()
    assert s.checkInclusion("ab", "eidbaooo") == True
    assert s.checkInclusion("ab", "eidboaoo") == False