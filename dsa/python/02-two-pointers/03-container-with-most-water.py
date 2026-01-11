# https://leetcode.com/problems/container-with-most-water/description/
# time complexity: O(n)
# space complexity: O(1)

class Solution:
    def maxArea(self, height: List[int]) -> int:
        l, r = 0, len(height)-1
        maxArea = 0

        while l < r:
            currHeight = min(height[l], height[r])
            maxArea = max(currHeight * (r-l), maxArea)
            
            if height[l] < height[r]:
                l += 1
            else:
                r -= 1
        
        return maxArea

if __name__ == "__main__":
    solution = Solution()

    assert solution.maxArea([1,8,6,2,5,4,8,3,7]) == 49
    assert solution.maxArea([1,1]) == 1
    assert solution.maxArea([]) == 0


