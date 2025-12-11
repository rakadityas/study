package linked_list

import "testing"

type RandomNode struct {
    Val    int
    Next   *RandomNode
    Random *RandomNode
}

func copyRandomList(head *RandomNode) *RandomNode {
    if head == nil { return nil }
    m := map[*RandomNode]*RandomNode{}
    for cur := head; cur != nil; cur = cur.Next {
        m[cur] = &RandomNode{Val: cur.Val}
    }
    for cur := head; cur != nil; cur = cur.Next {
        m[cur].Next = m[cur.Next]
        m[cur].Random = m[cur.Random]
    }
    return m[head]
}

func TestCopyRandomList(t *testing.T) {
    a := &RandomNode{Val:1}
    b := &RandomNode{Val:2}
    a.Next = b
    a.Random = b
    b.Random = a
    cp := copyRandomList(a)
    if cp.Val != 1 || cp.Random.Val != 2 || cp.Next.Random.Val != 1 { t.Fatalf("copy incorrect") }
}

