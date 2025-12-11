package linked_list

import (
    "reflect"
    "testing"
)

func reorderList(head *ListNode) *ListNode {
    if head == nil || head.Next == nil { return head }
    // find middle
    slow, fast := head, head
    for fast != nil && fast.Next != nil { slow = slow.Next; fast = fast.Next.Next }
    // reverse second half
    var prev *ListNode
    cur := slow.Next
    slow.Next = nil
    for cur != nil { nxt := cur.Next; cur.Next = prev; prev = cur; cur = nxt }
    // merge
    p1, p2 := head, prev
    for p2 != nil {
        t1, t2 := p1.Next, p2.Next
        p1.Next = p2
        p2.Next = t1
        p1, p2 = t1, t2
    }
    return head
}

func TestReorderList(t *testing.T) {
    got := toSlice(reorderList(buildList([]int{1,2,3,4})))
    want := []int{1,4,2,3}
    if !reflect.DeepEqual(got, want) { t.Fatalf("got %v want %v", got, want) }
}

