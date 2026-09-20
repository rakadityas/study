package tree

import "testing"

// diameterOfBinaryTree returns the longest path (in edges) between any two nodes.
// Approach: post-order DFS returning subtree height while tracking the best left+right sum seen
// time: O(n), space: O(n) — one visit per node; recursion stack is O(h), O(n) when skewed
func diameterOfBinaryTree(root *TreeNode) int {
    maxDiameter := 0

    var dfs func(node *TreeNode) int
    dfs = func(node *TreeNode) int {
        if node == nil { return 0 }

        leftVal := dfs(node.Left)
        rightVal := dfs(node.Right)

        if leftVal+rightVal > maxDiameter { maxDiameter = leftVal + rightVal }

        // height only needs the taller of the two sides
        if leftVal > rightVal { return 1 + leftVal }
        return 1 + rightVal
    }

    dfs(root)
    return maxDiameter
}

func TestDiameterOfBinaryTree(t *testing.T) {
    cases := []struct {
        in   []any
        want int
    }{
        {[]any{1, 2, 3, 4, 5}, 3},
        {[]any{1, 2}, 1},
        {[]any{1, 2, 3, 4, 5, 6, 7}, 4},
        {[]any{}, 0},
    }

    for _, c := range cases {
        if got := diameterOfBinaryTree(FromList(c.in)); got != c.want {
            t.Fatalf("diameterOfBinaryTree(%v) = %d, want %d", c.in, got, c.want)
        }
    }
}
