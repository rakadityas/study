package tree

import "testing"

// invertTree swaps every node's children recursively.
// Approach: recursive DFS — swap children at each node, then recurse into both subtrees
// time: O(n), space: O(n) — visits every node once; recursion stack is O(h), O(n) for a skewed tree
func invertTree(root *TreeNode) *TreeNode {
    if root == nil { return nil }

    root.Left, root.Right = root.Right, root.Left

    invertTree(root.Left)
    invertTree(root.Right)
    return root
}

// invertTreeIterative swaps every node's children using an explicit stack.
// Approach: iterative DFS with an explicit stack — avoids recursion stack overhead
// time: O(n), space: O(n) — stack holds at most O(h) nodes, O(n) worst case
func invertTreeIterative(root *TreeNode) *TreeNode {
    if root == nil { return nil }

    stack := []*TreeNode{root}
    for len(stack) > 0 {
        curr := stack[len(stack)-1]
        stack = stack[:len(stack)-1]
        curr.Left, curr.Right = curr.Right, curr.Left

        if curr.Left != nil { stack = append(stack, curr.Left) }
        if curr.Right != nil { stack = append(stack, curr.Right) }
    }

    return root
}

func TestInvertTree(t *testing.T) {
    cases := []struct {
        in   []any
        want []any
    }{
        {[]any{4, 2, 7, 1, 3, 6, 9}, []any{4, 7, 2, 9, 6, 3, 1}},
        {[]any{2, 1, 3}, []any{2, 3, 1}},
        {[]any{1}, []any{1}},
    }

    for _, c := range cases {
        if got := invertTree(FromList(c.in)).ToList(); !equalList(got, c.want) {
            t.Fatalf("invertTree(%v) = %v, want %v", c.in, got, c.want)
        }
        if got := invertTreeIterative(FromList(c.in)).ToList(); !equalList(got, c.want) {
            t.Fatalf("invertTreeIterative(%v) = %v, want %v", c.in, got, c.want)
        }
    }

    if invertTree(FromList([]any{})) != nil { t.Fatalf("empty tree should invert to nil") }
    if invertTreeIterative(FromList([]any{})) != nil { t.Fatalf("empty tree should invert to nil") }
}
