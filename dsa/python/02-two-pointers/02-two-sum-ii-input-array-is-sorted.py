# https://leetcode.com/problems/two-sum-ii-input-array-is-sorted/description/
# time complexity: O(n)
# space complexity: O(1)

class Solution:
    def twoSum(self, numbers: list[int], target:int) -> list[int]:
        l = 0
        r = len(numbers) - 1

        while l < r:
            sumNumbers = numbers[l] + numbers[r]
            if sumNumbers == target:
                return [l+1, r+1]
            elif sumNumbers > target:
                r -= 1
            else:
                l += 1
        
        return [0,0]

if __name__ == "__main__":
    solution = Solution()

    assert solution.twoSum([2,7,11,15], 9) == [1,2]
    assert solution.twoSum([2,3,4], 6) == [1,3]
    assert solution.twoSum([-1,0], -1) == [1,2]