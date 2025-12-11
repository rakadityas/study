package linked_list

import "testing"

func hasCycle(head *ListNode) bool {
    slow, fast := head, head
    for fast != nil && fast.Next != nil {
        slow = slow.Next
        fast = fast.Next.Next
        if slow == fast { return true }
    }
    return false
}

func TestHasCycle(t *testing.T) {
    a := &ListNode{Val:1}; b := &ListNode{Val:2}; c := &ListNode{Val:3}
    a.Next = b; b.Next = c; c.Next = b
    if !hasCycle(a) { t.Fatalf("expected cycle") }
}

