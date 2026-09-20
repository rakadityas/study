package tree

import "testing"

// isBalanced reports whether every node's subtrees differ in height by at most 1.
// Approach: post-order DFS returning (balanced, depth) so an imbalance short-circuits the whole walk
// time: O(n), space: O(n) — single pass; recursion stack is O(h), O(n) when skewed
func isBalanced(root *TreeNode) bool {
    balanced, _ := isBalancedDFS(root)
    return balanced
}

func isBalancedDFS(root *TreeNode) (bool, int) {
    if root == nil { return true, 0 }

    leftBalanced, leftDepth := isBalancedDFS(root.Left)
    if !leftBalanced { return false, 0 }

    rightBalanced, rightDepth := isBalancedDFS(root.Right)
    if !rightBalanced { return false, 0 }

    diff := leftDepth - rightDepth
    if diff < 0 { diff = -diff }
    if diff > 1 { return false, 0 }

    if leftDepth > rightDepth { return true, 1 + leftDepth }
    return true, 1 + rightDepth
}

func TestIsBalanced(t *testing.T) {
    cases := []struct {
        in   []any
        want bool
    }{
        {[]any{3, 9, 20, nil, nil, 15, 7}, true},
        {[]any{1, 2, 2, 3, 3, nil, nil, 4, 4}, false},
        {[]any{}, true},
    }

    for _, c := range cases {
        if got := isBalanced(FromList(c.in)); got != c.want {
            t.Fatalf("isBalanced(%v) = %v, want %v", c.in, got, c.want)
        }
    }
}
