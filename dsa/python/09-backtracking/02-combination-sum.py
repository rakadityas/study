# https://leetcode.com/problems/combination-sum/description/
# Time Complexity: Exponential (backtracking)
# Space Complexity: O(target / min(candidates))

from typing import List

class Solution:
    def combinationSum(self, candidates: List[int], target: int) -> List[List[int]]:
        res = []

        def backtracking(start: int, combination: List[int], currentSum: int):
            if currentSum == target:
                res.append(combination.copy())
                return
            
            if currentSum > target:
                return 
            
            for i in range(start, len(candidates)):
                combination.append(candidates[i])
                backtracking(i, combination, currentSum+candidates[i])
                combination.pop()
            
        backtracking(0, [], 0)
        return res

if __name__ == "__main__":
    s = Solution()
    assert sorted(s.combinationSum([2,3,6,7], 7)) == sorted([[2,2,3],[7]])
    assert sorted(s.combinationSum([2,3,5], 8)) == sorted([[2,2,2,2],[2,3,3],[3,5]])
    assert sorted(s.combinationSum([2], 1)) == sorted([])
        