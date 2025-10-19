# https://leetcode.com/problems/search-in-rotated-sorted-array/description/
# time complexity: O(log(n))
# space complexity: O(1)

class Solution:
    def search(self, nums: list[int], target:int) -> int:
        l = 0
        r = len(nums) - 1

        while l <= r:
            mid = (l+r)//2
            if nums[mid] == target:
                return mid
            
            if nums[l] <= nums[mid]:
                if target > nums[mid] or target < nums[l]:
                    l = mid + 1
                else:
                    r = mid - 1
            else:
                if target < nums[mid] or target > nums[r]:
                    r = mid - 1
                else:
                    l = mid + 1
        
        return -1

if __name__ == "__main__":
    solution = Solution()
    assert solution.search([4,5,6,7,0,1,2], 0) == 4
    assert solution.search([4,5,6,7,0,1,2], 3) == -1
    assert solution.search([1], 0) == -1
