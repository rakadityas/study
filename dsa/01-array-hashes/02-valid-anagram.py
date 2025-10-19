#  https://leetcode.com/problems/valid-anagram/description/
# time complexity: O(n)
# space complexity: O(n)

from collections import defaultdict

class Solution:
    def isAnagram(self, s: str, t: str) -> bool:
        mapS = defaultdict(int)
        mapT = defaultdict(int)

        for i in range(len(s)):
            mapS[s[i]] += 1
        
        for i in range(len(t)):
            mapT[t[i]] += 1
        
        return mapT == mapS

if __name__ == '__main__':
    solution = Solution()

    assert solution.isAnagram("anagram", "nagaram") == True
    assert solution.isAnagram("rat", "car") == False
    assert solution.isAnagram("ratt", "car") == False