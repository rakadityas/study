package tree

import "testing"

// kthSmallest returns the kth smallest value in a BST.
// Approach: in-order DFS collecting all values into a sorted slice, then index into it
// time: O(n), space: O(n) — visits every node and stores every value
func kthSmallest(root *TreeNode, k int) int {
    treeValues := []int{}

    var dfs func(node *TreeNode)
    dfs = func(node *TreeNode) {
        if node == nil { return }

        dfs(node.Left)
        treeValues = append(treeValues, node.Val)
        dfs(node.Right)
    }

    dfs(root)
    return treeValues[k-1]
}

// kthSmallestOptimal returns the same value without materialising the sorted slice.
// Approach: in-order DFS counting k down as nodes are visited; the node that drives k to 0 is the answer
// time: O(n) worst case, O(h+k) typical, space: O(h) — recursion stack only
func kthSmallestOptimal(root *TreeNode, k int) int {
    res := 0

    var dfs func(node *TreeNode)
    dfs = func(node *TreeNode) {
        if node == nil { return }

        dfs(node.Left)

        k--
        if k == 0 { res = node.Val }

        dfs(node.Right)
    }

    dfs(root)
    return res
}

// kthSmallestPruning stops descending as soon as the answer is known.
// Approach: same in-order count-down, with an early return once k reaches 0 so the right subtree is skipped
// time: O(h+k), space: O(h) — prunes every node visited after the kth
func kthSmallestPruning(root *TreeNode, k int) int {
    res := 0

    var dfs func(node *TreeNode)
    dfs = func(node *TreeNode) {
        if node == nil || k == 0 { return }

        dfs(node.Left)
        if k == 0 { return }

        k--
        if k == 0 {
            res = node.Val
            return
        }

        dfs(node.Right)
    }

    dfs(root)
    return res
}

func TestKthSmallest(t *testing.T) {
    cases := []struct {
        in   []any
        k    int
        want int
    }{
        {[]any{3, 1, 4, nil, 2}, 1, 1},
        {[]any{5, 3, 6, 2, 4, nil, nil, 1}, 3, 3},
        {[]any{5, 3, 6, 2, 4, nil, nil, 1}, 6, 6},
    }

    for _, c := range cases {
        if got := kthSmallest(FromList(c.in), c.k); got != c.want {
            t.Fatalf("kthSmallest(%v, %d) = %d, want %d", c.in, c.k, got, c.want)
        }
        if got := kthSmallestOptimal(FromList(c.in), c.k); got != c.want {
            t.Fatalf("kthSmallestOptimal(%v, %d) = %d, want %d", c.in, c.k, got, c.want)
        }
        if got := kthSmallestPruning(FromList(c.in), c.k); got != c.want {
            t.Fatalf("kthSmallestPruning(%v, %d) = %d, want %d", c.in, c.k, got, c.want)
        }
    }
}
