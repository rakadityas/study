# https://leetcode.com/problems/find-minimum-in-rotated-sorted-array/description/
# time complexity: O(log(n))
# space complexity: O(1)

class Solution:
    def findMin(self, nums: list[int]) -> int:
        l = 0
        r = len(nums) - 1

        while l <= r:
            mid = (l+r)//2
            if nums[l] > nums[r]:
                l = mid + 1
            else:
                r = mid - 1

        return nums[l]
            


if __name__ == "__main__":
    solution = Solution()

    assert solution.findMin([3,4,5,1,2]) == 1
    assert solution.findMin([4,5,6,7,0,1,2]) == 0
    assert solution.findMin([11,13,15,17]) == 11