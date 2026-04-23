package linked_list

import "testing"

// Approach: Floyd's cycle detection — slow moves 1 step, fast moves 2; they meet iff a cycle exists
// time: O(n), space: O(1) — only two pointers (no hashmap unlike Python's first variant)
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

