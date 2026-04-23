# https://leetcode.com/problems/permutations/description/
# time complexity: O(n * n!) | n from copy result and n! from backtracking
# space complexity: O(n)

from typing import List

# Approach: boolean array to track used indices — faster than hashing, direct index access
# time complexity: O(n·n!) — n! permutations, each costs O(n) to copy
# space complexity: O(n) — boolean record array + recursion stack depth O(n)
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

# Approach: set to track used indices — O(1) average lookup but set overhead vs plain array
# time complexity: O(n·n!)
# space complexity: O(n) — set size bounded by n + recursion stack
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

# Approach: dict to track used indices — functionally same as SolutionSet; dict has slightly more overhead than set
# time complexity: O(n·n!)
# space complexity: O(n)
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