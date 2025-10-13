# https://leetcode.com/problems/subsets-ii/

from typing import List
from collections import defaultdict

class SolutionOne:
    def subsetsWithDup(self, nums: List[int]) -> List[List[int]]:
        res = []
        nums.sort()

        def backtracking(i: int, lenNums: int, currentSet: List[int]):
            if i >= lenNums:
                res.append(currentSet.copy())
                return
            
            currentSet.append(nums[i])
            backtracking(i+1, lenNums, currentSet)

            currentSet.pop()
            while i+1 < lenNums and nums[i] == nums[i+1]:
                i += 1
            backtracking(i+1, lenNums, currentSet)
            return
        
        backtracking(0, len(nums), [])
        return res

class SolutionTwo: # less optimal, but easier to understand
    def subsetsWithDup(self, nums: List[int]) -> List[List[int]]:
        nums.sort()
        res = []
        seen = defaultdict(bool)

        def backtrack(start: int, path: List[int]):
            key = tuple(path)
            if key not in seen:
                seen[key] = True
                res.append(path.copy())

            for i in range(start, len(nums)):
                path.append(nums[i])
                backtrack(i + 1, path)
                path.pop()

        backtrack(0, [])
        return res  

if __name__ == "__main__":
    s = SolutionOne()
    assert sorted(s.subsetsWithDup([1,2,2])) == sorted([[],[1],[1,2],[1,2,2],[2],[2,2]])
    assert sorted(s.subsetsWithDup([0])) == sorted([[],[0]])
    assert sorted(s.subsetsWithDup([4,4,1,4])) == sorted([[],[1],[1,4],[1,4,4],[1,4,4,4],[4],[4,4],[4,4,4]])

    s = SolutionTwo()
    assert sorted(s.subsetsWithDup([1,2,2])) == sorted([[],[1],[1,2],[1,2,2],[2],[2,2]])
    assert sorted(s.subsetsWithDup([0])) == sorted([[],[0]])
    assert sorted(s.subsetsWithDup([4,4,1,4])) == sorted([[],[1],[1,4],[1,4,4],[1,4,4,4],[4],[4,4],[4,4,4]])