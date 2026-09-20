package tree

import (
    "math"
    "testing"
)

// isValidBST reports whether the tree satisfies the BST ordering globally, not just locally.
// Approach: DFS narrowing an open (min, max) range as it descends each branch
// time: O(n), space: O(n) — one visit per node; recursion stack is O(h), O(n) when skewed
func isValidBST(root *TreeNode) bool {
    return isValidBSTDFS(root, math.MinInt64, math.MaxInt64)
}

func isValidBSTDFS(root *TreeNode, minVal, maxVal int) bool {
    if root == nil { return true }

    if root.Val <= minVal || root.Val >= maxVal { return false }

    return isValidBSTDFS(root.Left, minVal, root.Val) &&
        isValidBSTDFS(root.Right, root.Val, maxVal)
}

func TestIsValidBST(t *testing.T) {
    cases := []struct {
        in   []any
        want bool
    }{
        {[]any{2, 1, 3}, true},
        {[]any{5, 1, 4, nil, nil, 3, 6}, false},
        {[]any{1}, true},
        {[]any{2, 2, 2}, false},
    }

    for _, c := range cases {
        if got := isValidBST(FromList(c.in)); got != c.want {
            t.Fatalf("isValidBST(%v) = %v, want %v", c.in, got, c.want)
        }
    }
}
