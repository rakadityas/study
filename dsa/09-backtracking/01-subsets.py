# https://leetcode.com/problems/subsets/

class Solution:
    def subsets(self, nums: list[int]) -> list[list[int]]:
        res = []

        def backtracking(i: int, lenNums: int, currentSet: list[int]):
            if i >= lenNums:
                res.append(currentSet.copy())
                return

            currentSet.append(nums[i])
            backtracking(i+1, lenNums, currentSet)

            currentSet.pop()
            backtracking(i+1, lenNums, currentSet)
            return
        
        backtracking(0, len(nums), [])
        return res

if __name__ == "__main__":
    s = Solution()
    assert sorted(s.subsets([1,2,3])) == sorted([[],[1],[2],[1,2],[3],[1,3],[2,3],[1,2,3]])
    assert sorted(s.subsets([0])) == sorted([[],[0]])
