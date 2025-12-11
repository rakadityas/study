package binary_search

import "testing"

func searchMatrix(matrix [][]int, target int) bool {
    if len(matrix) == 0 || len(matrix[0]) == 0 { return false }
    m, n := len(matrix), len(matrix[0])
    l, r := 0, m*n-1
    for l <= r {
        mid := l + (r-l)/2
        x := matrix[mid/n][mid%n]
        if x == target { return true }
        if x < target { l = mid + 1 } else { r = mid - 1 }
    }
    return false
}

func TestSearchMatrix(t *testing.T) {
    mat := [][]int{{1,3,5,7},{10,11,16,20},{23,30,34,60}}
    if !searchMatrix(mat, 3) { t.Fatalf("should find 3") }
    if searchMatrix(mat, 13) { t.Fatalf("should not find 13") }
}

