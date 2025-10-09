# https://leetcode.com/problems/permutations/

from typing import List

class SolutionList:
    def permute(self, nums: List[int]) -> List[List[int]]:
        self.res = []
        self.backtracking(nums, [], [False] * len(nums))
        return self.res

    def backtracking(self, nums: List[int], permutation: List[int], record: List[bool]):
        if len(nums) == len(permutation):
            self.res.append(permutation.copy())
            return
        
        for i in range(len(nums)):
            if record[i] == True:
                continue
            
            record[i] = True
            permutation.append(nums[i])
            self.backtracking(nums, permutation, record)
            permutation.pop()
            record[i] = False

class SolutionSet:
    def permute(self, nums: List[int]) -> List[List[int]]:
        self.res = []
        self.backtracking(nums, [],set())
        return self.res

    def backtracking(self, nums: List[int], permutation: List[int], setHistory: set()):
        if len(nums) == len(permutation):
            self.res.append(permutation.copy())
            return
        
        for i in range(len(nums)):
            if i in setHistory:
                continue
            
            setHistory.add(i)
            permutation.append(nums[i])
            self.backtracking(nums, permutation, setHistory)
            setHistory.remove(i)
            permutation.pop()

class SolutionDict:
    def permute(self, nums: List[int]) -> List[List[int]]:
        self.res = []
        self.backtracking(nums, [],{})
        return self.res

    def backtracking(self, nums: List[int], permutation: List[int], dictHistory: dict()):
        if len(nums) == len(permutation):
            self.res.append(permutation.copy())
            return
        
        for i in range(len(nums)):
            if i in dictHistory:
                continue
            
            dictHistory[i] = True
            permutation.append(nums[i])
            self.backtracking(nums, permutation, dictHistory)
            del dictHistory[i]
            permutation.pop()

if __name__ == "__main__":
    s = SolutionList()
    assert sorted(s.permute([1,2,3])) == sorted([[1,2,3],[1,3,2],[2,1,3],[2,3,1],[3,1,2],[3,2,1]])
    assert sorted(s.permute([0,1])) == sorted([[0,1],[1,0]])
    assert sorted(s.permute([1])) == sorted([[1]])

    s = SolutionSet()
    assert sorted(s.permute([1,2,3])) == sorted([[1,2,3],[1,3,2],[2,1,3],[2,3,1],[3,1,2],[3,2,1]])
    assert sorted(s.permute([0,1])) == sorted([[0,1],[1,0]])
    assert sorted(s.permute([1])) == sorted([[1]])

    s = SolutionDict()
    assert sorted(s.permute([1,2,3])) == sorted([[1,2,3],[1,3,2],[2,1,3],[2,3,1],[3,1,2],[3,2,1]])
    assert sorted(s.permute([0,1])) == sorted([[0,1],[1,0]])
    assert sorted(s.permute([1])) == sorted([[1]])