
class Solution:
    def search(self, nums: list[int], target: int) -> int:
        l = 0
        r = len(nums) - 1

        while l < r:
            mid = (l+r)//2
            if target == nums[mid]:
                return mid
            elif target > nums[mid]:
                l = mid+1
            else:
                r = mid-1
        
        return -1

if __name__ == "__main__":
    solution = Solution()

    assert solution.search([-1,0,3,5,9,12], 9) == 4
    assert solution.search([-1,0,3,5,9,12], 2) == -1