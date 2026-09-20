package tree

// TreeNode is the shared binary-tree node for this package.
// Test cases describe trees as level-order []any slices where nil marks a missing child,
// mirroring the LeetCode array format used by the Python bank.
type TreeNode struct {
    Val   int
    Left  *TreeNode
    Right *TreeNode
}

// FromList builds a tree from a level-order slice, skipping nil slots.
// time: O(n), space: O(n) — one queue entry per created node
func FromList(values []any) *TreeNode {
    if len(values) == 0 { return nil }
    root := &TreeNode{Val: values[0].(int)}
    queue := []*TreeNode{root}
    i := 1
    for len(queue) > 0 && i < len(values) {
        node := queue[0]
        queue = queue[1:]
        if i < len(values) && values[i] != nil {
            node.Left = &TreeNode{Val: values[i].(int)}
            queue = append(queue, node.Left)
        }
        i++
        if i < len(values) && values[i] != nil {
            node.Right = &TreeNode{Val: values[i].(int)}
            queue = append(queue, node.Right)
        }
        i++
    }
    return root
}

// ToList flattens a tree back to level-order form with trailing nils trimmed.
// time: O(n), space: O(n) — BFS queue plus the result slice
func (t *TreeNode) ToList() []any {
    result := []any{}
    queue := []*TreeNode{t}
    for len(queue) > 0 {
        node := queue[0]
        queue = queue[1:]
        if node != nil {
            result = append(result, node.Val)
            queue = append(queue, node.Left, node.Right)
        } else {
            result = append(result, nil)
        }
    }
    for len(result) > 0 && result[len(result)-1] == nil {
        result = result[:len(result)-1]
    }
    return result
}

// equalList compares two level-order slices element by element.
func equalList(a, b []any) bool {
    if len(a) != len(b) { return false }
    for i := range a {
        if a[i] != b[i] { return false }
    }
    return true
}
