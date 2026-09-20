package tree

import "testing"

// lowestCommonAncestor finds the LCA of p and q in a binary search tree.
// Approach: walk down from the root — the first node that splits p and q is the LCA
// time: O(n), space: O(1) — O(h) steps (log n when balanced), iterative so no recursion stack
func lowestCommonAncestor(root, p, q *TreeNode) *TreeNode {
    curr := root
    for curr != nil {
        if p.Val < curr.Val && q.Val < curr.Val {
            curr = curr.Left
        } else if p.Val > curr.Val && q.Val > curr.Val {
            curr = curr.Right
        } else {
            return curr
        }
    }
    return nil
}

func TestLowestCommonAncestor(t *testing.T) {
    cases := []struct {
        tree []any
        p, q int
        want int
    }{
        {[]any{6, 2, 8, 0, 4, 7, 9, nil, nil, 3, 5}, 2, 8, 6},
        {[]any{6, 2, 8, 0, 4, 7, 9, nil, nil, 3, 5}, 2, 4, 2},
        {[]any{2, 1}, 2, 1, 2},
    }

    for _, c := range cases {
        got := lowestCommonAncestor(FromList(c.tree), &TreeNode{Val: c.p}, &TreeNode{Val: c.q})
        if got == nil || got.Val != c.want {
            t.Fatalf("lowestCommonAncestor(%v, %d, %d) = %v, want %d", c.tree, c.p, c.q, got, c.want)
        }
    }
}
