package tree

import "testing"

// rightSideViewQueue returns the rightmost value at each depth.
// Approach: DFS visiting the right child first and bucketing every value by depth,
// then taking the head of each bucket
// time: O(n), space: O(n) — listTree stores all nodes grouped by depth
func rightSideViewQueue(root *TreeNode) []int {
    listTree := [][]int{}

    var dfs func(node *TreeNode, depth int)
    dfs = func(node *TreeNode, depth int) {
        if node == nil { return }

        if depth >= len(listTree) { listTree = append(listTree, []int{}) }

        listTree[depth] = append(listTree[depth], node.Val)

        dfs(node.Right, depth+1)
        dfs(node.Left, depth+1)
    }
    dfs(root, 0)

    res := []int{}
    for _, level := range listTree {
        res = append(res, level[0])
    }
    return res
}

// rightSideViewList returns the same view without the per-depth buckets.
// Approach: DFS right-first, appending only when depth == len(res) — the first visit at that level
// time: O(n), space: O(n) — result slice plus an O(h) recursion stack
func rightSideViewList(root *TreeNode) []int {
    res := []int{}

    var dfs func(node *TreeNode, depth int)
    dfs = func(node *TreeNode, depth int) {
        if node == nil { return }

        if depth == len(res) { res = append(res, node.Val) }

        dfs(node.Right, depth+1)
        dfs(node.Left, depth+1)
    }

    dfs(root, 0)
    return res
}

func equalInts(a, b []int) bool {
    if len(a) != len(b) { return false }
    for i := range a {
        if a[i] != b[i] { return false }
    }
    return true
}

func TestRightSideView(t *testing.T) {
    cases := []struct {
        in   []any
        want []int
    }{
        {[]any{1, 2, 3, nil, 5, nil, 4}, []int{1, 3, 4}},
        {[]any{}, []int{}},
        {[]any{1, 2, 3, 4, nil, nil, nil, 5}, []int{1, 3, 4, 5}},
        {[]any{1, nil, 3}, []int{1, 3}},
    }

    for _, c := range cases {
        if got := rightSideViewQueue(FromList(c.in)); !equalInts(got, c.want) {
            t.Fatalf("rightSideViewQueue(%v) = %v, want %v", c.in, got, c.want)
        }
        if got := rightSideViewList(FromList(c.in)); !equalInts(got, c.want) {
            t.Fatalf("rightSideViewList(%v) = %v, want %v", c.in, got, c.want)
        }
    }
}
