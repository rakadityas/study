# https://leetcode.com/problems/search-a-2d-matrix/description/
# time complexity: O(log(m*n))
# space complexity: O(1)

class Solution:
    def searchMatrix(self, matrix: list[list[int]], target: int) -> bool:
        for i in range(len(matrix)):
            l, r = 0, len(matrix[i])-1
            
            while l <= r:
                mid = (r+l)//2
                if target == matrix[i][mid]:
                    return True
                elif target > matrix[i][mid]:
                    l = mid + 1
                else:
                    r = mid - 1

        return False

if __name__ == "__main__":
    solution = Solution()

    assert solution.searchMatrix([[1,3,5,7],[10,11,16,20],[23,30,34,60]], 3) == True 
    assert solution.searchMatrix([[1,3,5,7],[10,11,16,20],[23,30,34,60]], 13) == False 
