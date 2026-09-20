package tree

import "testing"

// isSameTree reports whether two trees have identical structure and values.
// Approach: parallel DFS comparing both nodes at each step, short-circuiting on the first mismatch
// time: O(n), space: O(n) — recursion stack is O(h), O(n) for a skewed tree
func isSameTree(p *TreeNode, q *TreeNode) bool {
    if p == nil && q == nil { return true }
    if p == nil || q == nil { return false }

    if p.Val != q.Val { return false }

    if !isSameTree(p.Left, q.Left) { return false }

    return isSameTree(p.Right, q.Right)
}

func TestIsSameTree(t *testing.T) {
    cases := []struct {
        p, q []any
        want bool
    }{
        {[]any{1, 2, 3}, []any{1, 2, 3}, true},
        {[]any{1, 2}, []any{1, nil, 2}, false},
        {[]any{1, 2, 1}, []any{1, 1, 2}, false},
        {[]any{}, []any{}, true},
    }

    for _, c := range cases {
        if got := isSameTree(FromList(c.p), FromList(c.q)); got != c.want {
            t.Fatalf("isSameTree(%v, %v) = %v, want %v", c.p, c.q, got, c.want)
        }
    }
}
