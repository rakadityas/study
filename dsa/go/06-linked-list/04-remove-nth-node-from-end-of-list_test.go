package linked_list

import (
    "reflect"
    "testing"
)

// Approach: two-pointer gap — advance fast n steps ahead, then move both until fast.Next == nil; slow is at the node before target
// time: O(n), space: O(1) — single pass with a dummy head to handle edge case of removing the first node
func removeNthFromEnd(head *ListNode, n int) *ListNode {
    dummy := &ListNode{Next: head}
    fast, slow := dummy, dummy
    for i := 0; i < n; i++ { fast = fast.Next }
    for fast.Next != nil { fast = fast.Next; slow = slow.Next }
    slow.Next = slow.Next.Next
    return dummy.Next
}

func TestRemoveNthFromEnd(t *testing.T) {
    got := toSlice(removeNthFromEnd(buildList([]int{1,2,3,4,5}), 2))
    want := []int{1,2,3,5}
    if !reflect.DeepEqual(got, want) { t.Fatalf("got %v want %v", got, want) }
}

