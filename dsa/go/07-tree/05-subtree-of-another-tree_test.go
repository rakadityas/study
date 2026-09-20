package tree

import "testing"

// isSubtree reports whether subRoot appears as a whole subtree inside root.
// Approach: DFS over root, running a full structural comparison from each candidate node
// time: O(n*m), space: O(n+m) — every root node may trigger an m-node comparison
func isSubtree(root *TreeNode, subRoot *TreeNode) bool {
    if root == nil { return false }

    if checkSubroot(root, subRoot) { return true }

    return isSubtree(root.Left, subRoot) || isSubtree(root.Right, subRoot)
}

func checkSubroot(root *TreeNode, subRoot *TreeNode) bool {
    if root == nil && subRoot == nil { return true }

    if root == nil || subRoot == nil { return false }

    if root.Val != subRoot.Val { return false }

    return checkSubroot(root.Left, subRoot.Left) && checkSubroot(root.Right, subRoot.Right)
}

func TestIsSubtree(t *testing.T) {
    cases := []struct {
        root    []any
        subRoot []any
        want    bool
    }{
        {[]any{3, 4, 5, 1, 2}, []any{4, 1, 2}, true},
        {[]any{3, 4, 5, 1, 2, nil, nil, nil, nil, 0}, []any{4, 1, 2}, false},
    }

    for _, c := range cases {
        if got := isSubtree(FromList(c.root), FromList(c.subRoot)); got != c.want {
            t.Fatalf("isSubtree(%v, %v) = %v, want %v", c.root, c.subRoot, got, c.want)
        }
    }
}
