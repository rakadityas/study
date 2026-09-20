package tree

import "testing"

// levelOrder groups node values by depth, top level first.
// Approach: DFS carrying the depth, appending a new bucket the first time a depth is reached
// time: O(n), space: O(n) — result holds every value; recursion stack is O(h)
func levelOrder(root *TreeNode) [][]int {
    res := [][]int{}

    var dfs func(node *TreeNode, depth int)
    dfs = func(node *TreeNode, depth int) {
        if node == nil { return }

        // ensure the result has enough buckets for the current depth
        if depth >= len(res) { res = append(res, []int{}) }

        res[depth] = append(res[depth], node.Val)

        dfs(node.Left, depth+1)
        dfs(node.Right, depth+1)
    }

    dfs(root, 0)
    return res
}

func equal2D(a, b [][]int) bool {
    if len(a) != len(b) { return false }
    for i := range a {
        if len(a[i]) != len(b[i]) { return false }
        for j := range a[i] {
            if a[i][j] != b[i][j] { return false }
        }
    }
    return true
}

func TestLevelOrder(t *testing.T) {
    cases := []struct {
        in   []any
        want [][]int
    }{
        {[]any{3, 9, 20, nil, nil, 15, 7}, [][]int{{3}, {9, 20}, {15, 7}}},
        {[]any{1}, [][]int{{1}}},
        {[]any{}, [][]int{}},
    }

    for _, c := range cases {
        if got := levelOrder(FromList(c.in)); !equal2D(got, c.want) {
            t.Fatalf("levelOrder(%v) = %v, want %v", c.in, got, c.want)
        }
    }
}
