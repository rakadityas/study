package two_pointers

import "testing"

// maxArea calculates the maximum area of water container using two pointers.
func maxArea(height []int) int {
    i, j := 0, len(height)-1
    best := 0
    for i < j {
        h := height[i]
        if height[j] < h { h = height[j] }
        area := h * (j - i)
        if area > best { best = area }
        if height[i] < height[j] { i++ } else { j-- }
    }
    return best
}

func TestMaxArea(t *testing.T) {
    if got := maxArea([]int{1,8,6,2,5,4,8,3,7}); got != 49 {
        t.Fatalf("maxArea example expected 49, got %d", got)
    }
}

