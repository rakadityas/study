from typing import List

class Solution:
    def sortArray(self, nums: List[int]) -> List[int]:
        self.mergeSort(nums, 0, len(nums)-1)
        return nums
    
    def mergeSort(self, nums: List[int], left: int, right: int):     
        if left >= right:
            return
        
        mid = (right+left)//2
        self.mergeSort(nums, left, mid)
        self.mergeSort(nums, mid+1, right)

        self.merge(nums, left, mid, right)
    
    def merge(self, nums: List[int], left: int, mid: int, right: int):
        leftNums, rightNums = nums[left:mid+1], nums[mid+1:right+1]
        leftIdx, rightIdx = 0, 0

        numsIdx = left

        while leftIdx < len(leftNums) and rightIdx < len(rightNums):
            if leftNums[leftIdx] < rightNums[rightIdx]:
                nums[numsIdx] = leftNums[leftIdx]
                leftIdx += 1
            else:
                 nums[numsIdx] = rightNums[rightIdx]
                 rightIdx += 1
            numsIdx += 1
        
        while leftIdx < len(leftNums):
            nums[numsIdx] = leftNums[leftIdx]
            leftIdx += 1
            numsIdx += 1
        
        while rightIdx < len(rightNums):
            nums[numsIdx] = rightNums[rightIdx]
            rightIdx += 1
            numsIdx += 1
        return
    
if __name__ == "__main__":
    s = Solution()
    assert s.sortArray([5,2,3,1]) == [1,2,3,5]
    assert s.sortArray([5,1,1,2,0,0]) == [0,0,1,1,2,5]



        


