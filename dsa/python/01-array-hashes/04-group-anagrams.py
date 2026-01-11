# https://leetcode.com/problems/group-anagrams/description/
# time complexity: O(n * m)
# space complexity: O(n * m)

from collections import defaultdict

class Solution:
    def groupAnagrams(self, strs: list[str]) -> list[list[str]]:

        mapRespGroup = defaultdict(list)

        for i in range(len(strs)):
            arrWordCounter = [0] * 26
            for j in range(len(strs[i])):
                arrWordCounter[ord(strs[i][j])-ord('a')] += 1
                
            mapRespGroup[tuple(arrWordCounter)].append(strs[i])
        
        res = []
        for arrGroup in mapRespGroup.values():
            res.append(arrGroup)
        
        return res


if __name__ == "__main__":
    solution = Solution()

    assert sorted(solution.groupAnagrams(["eat","tea","tan","ate","nat","bat"])) == [["bat"],["eat", "tea", "ate"],["tan","nat"]]
    assert solution.groupAnagrams([""]) == [[""]]
    assert solution.groupAnagrams(["a"]) == [["a"]]