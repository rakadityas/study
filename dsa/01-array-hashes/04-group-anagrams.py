from collections import defaultdict

class Solution:
    def groupAnagrams(self, strs: list[str]) -> list[list[str]]:

        mapRespGroup = defaultdict(list)

        for i in range(len(strs)):
            arrWordCounter = [0] * 26
            for j in range(len(strs[i])):
                arrWordCounter[ord('a')-ord(strs[i][j])] += 1
                
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