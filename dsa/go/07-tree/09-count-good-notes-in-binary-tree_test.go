package tree

import "testing"

// goodNodes counts nodes whose value is at least as large as every value on the path from the root.
// Approach: DFS carrying the maximum seen so far down each branch
// time: O(n), space: O(n) — each node is visited once; recursion stack is O(h), O(n) when skewed
func goodNodes(root *TreeNode) int {
    if root == nil { return 0 }

    count := 0

    var dfs func(node *TreeNode, maxVal int)
    dfs = func(node *TreeNode, maxVal int) {
        if node == nil { return }

        if node.Val >= maxVal {
            count++
            maxVal = node.Val
        }

        dfs(node.Left, maxVal)
        dfs(node.Right, maxVal)
    }

    dfs(root, root.Val)
    return count
}

func TestGoodNodes(t *testing.T) {
    cases := []struct {
        in   []any
        want int
    }{
        {[]any{3, 1, 4, 3, nil, 1, 5}, 4},
        {[]any{3, 3, nil, 4, 2}, 3},
        {[]any{1}, 1},
        {[]any{-1, 5, -2, 4, 4, 2, -2, nil, nil, -4, nil, -2, 3, nil, -2, 0, nil, -1, nil, -3, nil, -4, -3, 3, nil, nil, nil, nil, nil, nil, nil, 3, -3}, 5},
    }

    for _, c := range cases {
        if got := goodNodes(FromList(c.in)); got != c.want {
            t.Fatalf("goodNodes(%v) = %d, want %d", c.in, got, c.want)
        }
    }
}
