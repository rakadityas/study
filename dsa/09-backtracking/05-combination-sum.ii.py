# https://leetcode.com/problems/combination-sum-ii/description/

from typing import List

class SolutionOne:
    def combinationSum2(self, candidates: List[int], target: int) -> List[List[int]]:
        candidates.sort(reverse=True)
        res = []
        mapHistory = {}

        def backtracking(combination: List[int], currentSum: int, start: int):
            if currentSum == target:
                combinationTuple = tuple(combination)
                if combinationTuple not in mapHistory:
                    res.append(combination.copy())
                    mapHistory[combinationTuple] = True
                return
                
            elif currentSum > target:
                return
            
            for i in range(start, len(candidates)):
                combination.append(candidates[i])
                backtracking(combination, currentSum+candidates[i], i+1)
                combination.pop()
            
            return
        
        backtracking([], 0, 0)

        return res

class SolutionTwo:
    def combinationSum2(self, candidates: List[int], target: int) -> List[List[int]]:
        res = []
        candidates.sort()

        def backtracking(path: List[int], currentSum: int, start: int):
            if currentSum == target:
                res.append(path.copy())
                return
            
            elif currentSum > target:
                return
            
            prev = None
            for i in range(start, len(candidates)):
                if candidates[i] == prev:
                    continue

                path.append(candidates[i])
                backtracking(path, candidates[i]+currentSum, i+1)
                path.pop()

                prev = candidates[i]
            return

        backtracking([], 0, 0)
        return res




if __name__ == "__main__":
    s = SolutionTwo()
    assert s.combinationSum2([10,1,2,7,6,1,5], 8).sort() == [[1,1,6],[1,2,5],[1,7],[2,6]].sort()
    # assert s.combinationSum2([2,5,2,1,2],5).sort() == [[1,2,2],[5]].sort()
    # assert s.combinationSum2([1,1,1,1,1,1],5).sort() == [[1,2,2],[5]].sort()
    # assert s.combinationSum2([1,1,1,1,1,1,1,1,1,1,1,1,1,1,1,1,1,1,1,1,1,1,1,1,1,1,1,1,1,1,1,1,1,1,1,1,1,1,1,1,1,1,1,1,1,1,1,1,1,1,1,1,1,1,1,1,1,1,1,1,1,1,1,1,1,1,1,1,1,1,1,1,1,1,1,1,1,1,1,1,1,1,1,1,1,1,1,1,1,1,1,1,1,1,1,1,1,1,1,1], 30) == 10 ## super slow
            
            



        
