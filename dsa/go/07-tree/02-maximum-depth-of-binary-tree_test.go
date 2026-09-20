package tree

import "testing"

// maxDepth returns the number of nodes on the longest root-to-leaf path.
// Approach: recursive DFS carrying the running depth down the tree
// time: O(n), space: O(n) — recursion stack is O(h), O(n) for a skewed tree
func maxDepth(root *TreeNode) int {
    return maxDepthDFS(root, 0)
}

func maxDepthDFS(root *TreeNode, counter int) int {
    if root == nil { return counter }

    left := maxDepthDFS(root.Left, counter+1)
    right := maxDepthDFS(root.Right, counter+1)
    if left > right { return left }
    return right
}

// maxDepthIterative returns the same depth without recursion.
// Approach: iterative DFS with an explicit stack storing (node, depth) pairs
// time: O(n), space: O(n) — explicit stack replaces the recursion stack
func maxDepthIterative(root *TreeNode) int {
    if root == nil { return 0 }

    type frame struct {
        node  *TreeNode
        depth int
    }

    best := 0
    stack := []frame{{root, 1}}

    for len(stack) > 0 {
        f := stack[len(stack)-1]
        stack = stack[:len(stack)-1]

        if f.depth > best { best = f.depth }

        if f.node.Left != nil { stack = append(stack, frame{f.node.Left, f.depth + 1}) }
        if f.node.Right != nil { stack = append(stack, frame{f.node.Right, f.depth + 1}) }
    }

    return best
}

func TestMaxDepth(t *testing.T) {
    cases := []struct {
        in   []any
        want int
    }{
        {[]any{3, 9, 20, nil, nil, 15, 7}, 3},
        {[]any{1, nil, 2}, 2},
        {[]any{}, 0},
    }

    for _, c := range cases {
        if got := maxDepth(FromList(c.in)); got != c.want {
            t.Fatalf("maxDepth(%v) = %d, want %d", c.in, got, c.want)
        }
        if got := maxDepthIterative(FromList(c.in)); got != c.want {
            t.Fatalf("maxDepthIterative(%v) = %d, want %d", c.in, got, c.want)
        }
    }
}
