package graph

import "testing"

type GraphNode struct {
    Val int
    Neighbors []*GraphNode
}

func cloneGraph(node *GraphNode) *GraphNode {
    if node == nil { return nil }
    m := map[*GraphNode]*GraphNode{}
    var dfs func(n *GraphNode) *GraphNode
    dfs = func(n *GraphNode) *GraphNode {
        if n == nil { return nil }
        if cp, ok := m[n]; ok { return cp }
        cp := &GraphNode{Val: n.Val}
        m[n] = cp
        for _, nb := range n.Neighbors { cp.Neighbors = append(cp.Neighbors, dfs(nb)) }
        return cp
    }
    return dfs(node)
}

func TestCloneGraph(t *testing.T) {
    a := &GraphNode{Val:1}; b := &GraphNode{Val:2}
    a.Neighbors = []*GraphNode{b}
    b.Neighbors = []*GraphNode{a}
    cp := cloneGraph(a)
    if cp == a || cp.Neighbors[0] == b { t.Fatalf("deep copy expected") }
    if cp.Val != 1 || cp.Neighbors[0].Val != 2 { t.Fatalf("values mismatch") }
}

