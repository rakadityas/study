# https://leetcode.com/problems/container-with-most-water/description/
# time complexity: O(n)
# space complexity: O(1)

class Solution:
    def maxArea(self, height: list[int]) -> int:

        l = 0
        r = len(height)-1

        highestNum = 0

        while l < r:
            baseNum = height[r]
            if height[l] < height[r]:
                baseNum = height[l]
            
            value = baseNum * (r-l)

            if highestNum < value:
                highestNum = value
            
            if height[l] > height[r]:
                r -= 1
            else:
                l += 1
            
        return highestNum

if __name__ == "__main__":
    solution = Solution()

    assert solution.maxArea([1,8,6,2,5,4,8,3,7]) == 49
    assert solution.maxArea([1,1]) == 1
    assert solution.maxArea([]) == 0


