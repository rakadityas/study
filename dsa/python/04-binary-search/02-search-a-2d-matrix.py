# https://leetcode.com/problems/search-a-2d-matrix/description/
# time complexity: O(log(m*n))
# space complexity: O(1)

class Solution:
    def searchMatrix(self, matrix: List[List[int]], target: int) -> bool:
        if not matrix or not matrix[0]:
            return False

        m, n = len(matrix), len(matrix[0])
        l, r = 0, m * n - 1

        while l <= r:
            mid = (l + r) // 2
            val = matrix[mid // n][mid % n]

            if val == target:
                return True
            elif val < target:
                l = mid + 1
            else:
                r = mid - 1

        return False

if __name__ == "__main__":
    solution = Solution()

    assert solution.searchMatrix([[1,3,5,7],[10,11,16,20],[23,30,34,60]], 3) == True 
    assert solution.searchMatrix([[1,3,5,7],[10,11,16,20],[23,30,34,60]], 13) == False 
